// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use tauri_plugin_log::{Target, TargetKind, RotationStrategy, TimezoneStrategy};
use std::path::PathBuf;

fn main() {
    // Get project directory (src-tauri folder)
    // Menggunakan CARGO_MANIFEST_DIR untuk mendapatkan path ke folder src-tauri
    let project_dir = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    
    // Create logs directory path relative to project
    // Logs akan disimpan di frontend/src-tauri/logs/
    let logs_dir = project_dir.join("logs");
    
    tauri::Builder::default()
        .plugin(
            tauri_plugin_log::Builder::new()
                // Log ke terminal (stdout)
                .target(Target::new(TargetKind::Stdout))
                // Log ke webview console
                .target(Target::new(TargetKind::Webview))
                // Log ke file di folder logs project untuk troubleshooting
                .target(
                    Target::new(TargetKind::Folder {
                        path: logs_dir,
                        file_name: Some("warehouse-management".to_string()),
                    })
                )
                // Set maximum log level (Info = info, warn, error)
                .level(log::LevelFilter::Info)
                // Rotate log file instead of discarding
                .rotation_strategy(RotationStrategy::KeepAll)
                // Maximum file size: 10MB
                .max_file_size(10_000_000)
                // Use local timezone instead of UTC
                .timezone_strategy(TimezoneStrategy::UseLocal)
                .build(),
        )
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

