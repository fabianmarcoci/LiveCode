#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
mod window;

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            window::adjust_window_size(app);
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("App start error.");
}
