import "./App.css";
import { useEffect, useState } from "react";
import { getCurrentWindow } from "@tauri-apps/api/window";

const appWindow = getCurrentWindow();

function App() {
  const [isMaximized, setIsMaximized] = useState(false);

  useEffect(() => {
    const appWindow = getCurrentWindow();

    const updateMaximizedState = async () => {
      const max = await appWindow.isMaximized();
      setIsMaximized(max);
    };

    updateMaximizedState();

    const unlistenPromise = appWindow.listen("tauri://resize", () => {
      updateMaximizedState();
    });

    return () => {
      unlistenPromise.then((unlisten) => unlisten());
    };
  }, []);

  return (
    <>
      <div className="titlebar" data-tauri-drag-region>
        <img src="/white-logo.svg" className="titlebar-icon" alt="logo" />

        <div className="titlebar-buttons">
          <button className="titlebar-btn" onClick={() => appWindow.minimize()}>
            ⎯
          </button>

          <button
            className="titlebar-btn"
            onClick={async () => {
              if (await appWindow.isMaximized()) {
                await appWindow.unmaximize();
                setIsMaximized(false);
              } else {
                await appWindow.maximize();
                setIsMaximized(true);
              }
            }}
          >
            {isMaximized ? "❐" : "☐"}
          </button>

          <button
            className="titlebar-btn close-btn"
            onClick={() => appWindow.close()}
          >
            ✕
          </button>
        </div>
      </div>

      <main className="container"></main>
    </>
  );
}

export default App;
