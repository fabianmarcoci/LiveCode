#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]
mod commands;

use dotenvy::dotenv;
use tauri::Manager;

use commands::{register_user, check_email_available, check_username_available};
use live_code_lib::{adjust_window_size, init_db};

fn main() {
    dotenv().ok();
    tauri::Builder::default()
        .setup(|app| {
            tauri::async_runtime::block_on(async {
                let pool = init_db().await;
                app.manage(pool);
            });

            adjust_window_size(app);

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            register_user,
            check_email_available,
            check_username_available
        ])
        .run(tauri::generate_context!())
        .expect("App start error.");
}
