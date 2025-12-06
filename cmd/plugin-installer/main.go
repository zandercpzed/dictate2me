package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// ObsidianConfig represents the structure of obsidian.json
type ObsidianConfig struct {
	Vaults map[string]struct {
		Path string `json:"path"`
		Ts   int64  `json:"ts"`
		Open bool   `json:"open"`
	} `json:"vaults"`
}

func main() {
	fmt.Println("🎤 dictate2me - Plugin Installer")
	fmt.Println("===============================")
	fmt.Println()

	// 1. Locate Obsidian config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fatal("Could not find user home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, "Library", "Application Support", "obsidian", "obsidian.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try Linux path just in case
		configPath = filepath.Join(homeDir, ".config", "obsidian", "obsidian.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fatal("Obsidian configuration not found. Is Obsidian installed?")
		}
	}

	// 2. Read config to find vaults
	data, err := os.ReadFile(configPath)
	if err != nil {
		fatal("Could not read Obsidian config: %v", err)
	}

	var config ObsidianConfig
	if err := json.Unmarshal(data, &config); err != nil {
		fatal("Could not parse Obsidian config: %v", err)
	}

	if len(config.Vaults) == 0 {
		fatal("No Obsidian vaults found in configuration.")
	}

	// 3. List vaults and ask user
	fmt.Println("Found the following vaults:")
	var paths []string
	i := 1
	for _, v := range config.Vaults {
		fmt.Printf("[%d] %s\n", i, v.Path)
		paths = append(paths, v.Path)
		i++
	}
	fmt.Println()
	fmt.Print("Select target vault (number): ")

	var selection int
	_, err = fmt.Scanf("%d", &selection)
	if err != nil || selection < 1 || selection > len(paths) {
		fatal("Invalid selection.")
	}

	targetVault := paths[selection-1]
	fmt.Printf("\nTarget vault: %s\n", targetVault)

	// 4. Build Plugin
	pluginDir := "./plugins/obsidian-dictate2me"
	fmt.Println("\nBuilding plugin...")

	// Check if npm is installed
	if _, err := exec.LookPath("npm"); err != nil {
		fatal("npm is not installed. Please install Node.js and npm.")
	}

	// npm install
	fmt.Println("Running npm install...")
	cmd := exec.Command("npm", "install")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("npm install failed: %v", err)
	}

	// npm run build
	fmt.Println("Running npm run build...")
	cmd = exec.Command("npm", "run", "build")
	cmd.Dir = pluginDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("npm run build failed: %v", err)
	}

	// 5. Install (Copy files)
	targetPluginDir := filepath.Join(targetVault, ".obsidian", "plugins", "obsidian-dictate2me")
	fmt.Printf("\nInstalling to: %s\n", targetPluginDir)

	if err := os.MkdirAll(targetPluginDir, 0755); err != nil {
		fatal("Could not create plugin directory: %v", err)
	}

	filesToCopy := []string{"main.js", "manifest.json", "styles.css"}
	for _, f := range filesToCopy {
		src := filepath.Join(pluginDir, f)
		dst := filepath.Join(targetPluginDir, f)
		if err := copyFile(src, dst); err != nil {
			fatal("Failed to copy %s: %v", f, err)
		}
		fmt.Printf("✓ Copied %s\n", f)
	}

	// 6. Configures Hot Reload (optional, creates .hotreload file)
	// This signals the Hot Reload plugin (if installed) to watch this folder
	os.WriteFile(filepath.Join(targetPluginDir, ".hotreload"), []byte(""), 0644)

	fmt.Println("\n✅ Plugin installed successfully!")
	// 7. Open Obsidian
	fmt.Println("\n🚀 Opening Obsidian project...")
	openCmd := exec.Command("open", "obsidian://open?path="+targetVault)
	if err := openCmd.Run(); err != nil {
		fmt.Printf("⚠️ Could not open Obsidian automatically: %v\n", err)
	} else {
		fmt.Println("✓ Obsidian opened!")
	}

	fmt.Println("--------------------------------")
	fmt.Println("Next steps in Obsidian:")
	fmt.Println("1. Go to Settings > Community Plugins")
	fmt.Println("2. Enable 'dictate2me'")
	fmt.Println("3. Enjoy!")
}

func fatal(format string, args ...interface{}) {
	fmt.Printf("\n❌ Error: "+format+"\n", args...)
	os.Exit(1)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
