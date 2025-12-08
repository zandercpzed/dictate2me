import { SimulatedCorrector } from './correction';

console.log('[TEST][CORRECTION]', { status: 'start' });
const corrector = new SimulatedCorrector();
const result = corrector.correct('texto de teste');
console.log('[TEST][CORRECTION]', { result });
