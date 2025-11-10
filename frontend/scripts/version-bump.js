#!/usr/bin/env node

/**
 * Script untuk auto-increment version dan create git tag
 * Usage:
 *   node scripts/version-bump.js patch   # 1.1.0 -> 1.1.1
 *   node scripts/version-bump.js minor   # 1.1.0 -> 1.2.0
 *   node scripts/version-bump.js major   # 1.1.0 -> 2.0.0
 */

import { readFileSync, writeFileSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';
import { execSync } from 'child_process';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const rootDir = join(__dirname, '..');

// Parse version string (e.g., "1.2.3") to object
function parseVersion(version) {
  const parts = version.split('.').map(Number);
  return {
    major: parts[0],
    minor: parts[1],
    patch: parts[2],
    toString() {
      return `${this.major}.${this.minor}.${this.patch}`;
    }
  };
}

// Bump version based on type
function bumpVersion(currentVersion, type) {
  const version = parseVersion(currentVersion);
  
  switch (type) {
    case 'major':
      version.major++;
      version.minor = 0;
      version.patch = 0;
      break;
    case 'minor':
      version.minor++;
      version.patch = 0;
      break;
    case 'patch':
      version.patch++;
      break;
    default:
      throw new Error(`Invalid version type: ${type}. Use 'major', 'minor', or 'patch'`);
  }
  
  return version.toString();
}

// Update version in JSON file
function updateJsonFile(filePath, newVersion) {
  const content = JSON.parse(readFileSync(filePath, 'utf-8'));
  content.version = newVersion;
  writeFileSync(filePath, JSON.stringify(content, null, 2) + '\n');
  console.log(`✓ Updated ${filePath} to version ${newVersion}`);
}

// Update version in Cargo.toml
function updateCargoToml(filePath, newVersion) {
  let content = readFileSync(filePath, 'utf-8');
  // Replace version = "x.y.z" with new version
  content = content.replace(/^version = "[\d.]+"/m, `version = "${newVersion}"`);
  writeFileSync(filePath, content);
  console.log(`✓ Updated ${filePath} to version ${newVersion}`);
}

// Create git tag
function createGitTag(version) {
  try {
    const tagName = `v${version}`;
    // Check if tag already exists
    try {
      execSync(`git rev-parse ${tagName}`, { stdio: 'ignore' });
      console.log(`⚠ Tag ${tagName} already exists, skipping tag creation`);
      return;
    } catch (e) {
      // Tag doesn't exist, create it
    }
    
    execSync(`git tag -a ${tagName} -m "Release version ${version}"`, { stdio: 'inherit' });
    console.log(`✓ Created git tag: ${tagName}`);
  } catch (error) {
    console.error(`⚠ Failed to create git tag: ${error.message}`);
    console.log('  You can create it manually: git tag -a v' + version + ' -m "Release version ' + version + '"');
  }
}

// Main function
function main() {
  const versionType = process.argv[2] || 'patch';
  
  if (!['major', 'minor', 'patch'].includes(versionType)) {
    console.error('Error: Invalid version type. Use "major", "minor", or "patch"');
    process.exit(1);
  }
  
  // Read current version from tauri.conf.json (source of truth)
  const tauriConfigPath = join(rootDir, 'src-tauri', 'tauri.conf.json');
  const tauriConfig = JSON.parse(readFileSync(tauriConfigPath, 'utf-8'));
  const currentVersion = tauriConfig.version;
  
  console.log(`Current version: ${currentVersion}`);
  console.log(`Bumping ${versionType} version...`);
  
  // Calculate new version
  const newVersion = bumpVersion(currentVersion, versionType);
  console.log(`New version: ${newVersion}\n`);
  
  // Update all version files
  console.log('Updating version files...');
  
  // 1. tauri.conf.json (source of truth)
  updateJsonFile(tauriConfigPath, newVersion);
  
  // 2. Cargo.toml
  const cargoTomlPath = join(rootDir, 'src-tauri', 'Cargo.toml');
  updateCargoToml(cargoTomlPath, newVersion);
  
  // 3. package.json (optional, for consistency)
  const packageJsonPath = join(rootDir, 'package.json');
  try {
    updateJsonFile(packageJsonPath, newVersion);
  } catch (error) {
    console.log(`⚠ Could not update package.json: ${error.message}`);
  }
  
  // Create git tag
  console.log('\nCreating git tag...');
  createGitTag(newVersion);
  
  console.log(`\n✅ Version bumped to ${newVersion} successfully!`);
  console.log(`\nNext steps:`);
  console.log(`  1. Review the changes: git diff`);
  console.log(`  2. Commit the changes: git add . && git commit -m "Bump version to ${newVersion}"`);
  console.log(`  3. Push with tags: git push && git push --tags`);
}

main();

