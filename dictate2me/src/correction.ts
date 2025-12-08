// Módulo de correção de texto (simulado para MVP)
// Interface e implementação fake para testes e integração

export interface Corrector {
  correct(text: string): string;
}

export class SimulatedCorrector implements Corrector {
  // Simula correção: adiciona pontuação
  correct(text: string): string {
    const corrected = text + '.';
    console.log('[CORRECTION]', { event: 'correct', original: text, corrected });
    return corrected;
  }
}

// Correção real via API local (HTTP)
export class ApiCorrector implements Corrector {
  async correct(text: string): Promise<string> {
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 4000);
      const response = await fetch('http://localhost:8080/api/v1/correct', {
        method: 'POST',
        body: JSON.stringify({ text }),
        headers: { 'Content-Type': 'application/json' },
        signal: controller.signal
      });
      clearTimeout(timeout);
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const { corrected } = await response.json();
      console.log('[API][CORRECT]', { corrected });
      return corrected || '';
    } catch (err) {
      console.log('[API][CORRECT][ERRO]', { err });
      return '';
    }
  }
}
