import "./App.css";
import { useEffect, useState } from "react";
import { getCurrentWindow } from "@tauri-apps/api/window";

const appWindow = getCurrentWindow();

function App() {
  const [isMaximized, setIsMaximized] = useState(false);
  const [showViewMenu, setShowViewMenu] = useState(false);
  const [showHiddenFiles, setShowHiddenFiles] = useState(false);
  const [showOptionsMenu, setShowOptionsMenu] = useState(false);
  const [lightMode, setLightMode] = useState(false);

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

  useEffect(() => {
    if (lightMode) {
      document.body.classList.add("light");
    } else {
      document.body.classList.remove("light");
    }
  }, [lightMode]);

  return (
    <>
      <div className="titlebar" data-tauri-drag-region>
        <div className="left-controls" data-tauri-drag-region={false}>
          <img
            src={lightMode ? "/black-logo.svg" : "/white-logo.svg"}
            className="titlebar-icon"
            alt="logo"
          />

          <div className="menu-group">
            <div className="menu">
              <button
                className={`menu-btn ${showViewMenu ? "active" : ""}`}
                onClick={() => {
                  setShowOptionsMenu(false);
                  setShowViewMenu((prev) => !prev);
                }}
              >
                View
              </button>

              {showViewMenu && (
                <div className="menu-dropdown">
                  <div
                    className="menu-item"
                    onClick={() => setShowHiddenFiles(!showHiddenFiles)}
                  >
                    <span className="menu-text">Show hidden files</span>

                    <div
                      className={`toggle-switch ${
                        showHiddenFiles ? "on" : "off"
                      }`}
                    >
                      <div className="toggle-knob" />
                    </div>
                  </div>
                </div>
              )}
            </div>

            <div className="menu">
              <button
                className={`menu-btn ${showOptionsMenu ? "active" : ""}`}
                onClick={() => {
                  setShowViewMenu(false);
                  setShowOptionsMenu((prev) => !prev);
                }}
              >
                Options
              </button>

              {showOptionsMenu && (
                <div className="menu-dropdown">
                  <div
                    className="menu-item"
                    onClick={() => setLightMode(!lightMode)}
                  >
                    <span className="menu-text">Light mode</span>

                    <div
                      className={`toggle-switch ${lightMode ? "on" : "off"}`}
                    >
                      <div className="toggle-knob" />
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

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
