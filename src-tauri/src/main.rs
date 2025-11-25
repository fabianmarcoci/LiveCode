#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
use dotenvy::dotenv;
use tauri::Manager;

mod database;
mod window;

use database::init_db;
fn main() {
    dotenv().ok();
    tauri::Builder::default()
        .setup(|app| {
            tauri::async_runtime::block_on(async {
                let pool = init_db().await;

                app.manage(pool);
            });

            window::adjust_window_size(app);

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("App start error.");
}
