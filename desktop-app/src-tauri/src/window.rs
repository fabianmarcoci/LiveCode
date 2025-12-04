use tauri::{App, Manager, PhysicalSize, Size};

pub fn adjust_window_size(app: &App) {
    let monitor = match app.primary_monitor() {
        Ok(Some(m)) => m,
        Ok(None) => {
            eprintln!("Warning: No primary monitor detected, skipping window adjustment");
            return;
        }
        Err(e) => {
            eprintln!("Error: Failed to get primary monitor: {:?}", e);
            return;
        }
    };
    let size = monitor.size();

    if size.width < 1200 || size.height < 800 {
        if let Some(window) = app.get_webview_window("main") {
            let new_width = size.width.saturating_sub(50);
            let new_height = size.height.saturating_sub(50);

            window
                .set_size(Size::Physical(PhysicalSize {
                    width: new_width,
                    height: new_height,
                }))
                .ok();
        }
    }
}
