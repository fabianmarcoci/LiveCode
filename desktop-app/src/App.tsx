import "./App.css";
import { useEffect, useState, useRef } from "react";
import { getCurrentWindow } from "@tauri-apps/api/window";
import { invoke } from "@tauri-apps/api/core";
import RegisterForm from "./components/auth/RegisterForm";
import LoginForm from "./components/auth/LoginForm";

const appWindow = getCurrentWindow();

interface User {
  username: string;
  email: string;
}

function App() {
  const [isMaximized, setIsMaximized] = useState(false);
  const [showViewMenu, setShowViewMenu] = useState(false);
  const [showHiddenFiles, setShowHiddenFiles] = useState(false);
  const [showOptionsMenu, setShowOptionsMenu] = useState(false);
  const [lightMode, setLightMode] = useState(false);
  const [showAuthMenu, setShowAuthMenu] = useState(false);
  const [authPanelType, setAuthPanelType] = useState<
    "register" | "login" | null
  >(null);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [showUserMenu, setShowUserMenu] = useState(false);

  const viewRef = useRef<HTMLDivElement | null>(null);
  const optionsRef = useRef<HTMLDivElement | null>(null);
  const authRef = useRef<HTMLDivElement | null>(null);
  const userMenuRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    async function checkAuth() {
      try {
        const token = await invoke<string | null>("get_access_token");

        if (token) {
          const userData = await invoke<User | null>("get_user_profile", {
            token,
          });

          if (userData) {
            setCurrentUser({
              username: userData.username,
              email: userData.email,
            });
          }
        }
      } catch (error) {
        console.error("Auth check failed:", error);
      }
    }

    checkAuth();
  }, []);

  // Maximize button
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

  // Light mode
  useEffect(() => {
    if (lightMode) {
      document.body.classList.add("light");
    } else {
      document.body.classList.remove("light");
    }
  }, [lightMode]);

  // Click outside dropdowns
  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      const target = e.target as Node;

      const insideView = viewRef.current?.contains(target);
      const insideOptions = optionsRef.current?.contains(target);
      const insideAuth = authRef.current?.contains(target);
      const insideUserMenu = userMenuRef.current?.contains(target);

      const clickedMenuButton = (target as HTMLElement).closest(
        ".menu-btn, .auth-btn, .user-btn"
      );

      if (
        insideView ||
        insideOptions ||
        insideAuth ||
        insideUserMenu ||
        clickedMenuButton
      ) {
        return;
      }

      setShowViewMenu(false);
      setShowOptionsMenu(false);
      setShowAuthMenu(false);
      setShowUserMenu(false);
    }

    window.addEventListener("click", handleClickOutside);

    return () => window.removeEventListener("click", handleClickOutside);
  }, []);

  // Register / Login panels
  useEffect(() => {
    if (!authPanelType) return;

    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") {
        setAuthPanelType(null);
      }
    }

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [authPanelType]);

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
                  setShowAuthMenu(false);
                  setShowViewMenu((prev) => !prev);
                }}
              >
                View
              </button>

              {showViewMenu && (
                <div className="menu-dropdown" ref={viewRef}>
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
                  setShowAuthMenu(false);
                  setShowOptionsMenu((prev) => !prev);
                }}
              >
                Options
              </button>

              {showOptionsMenu && (
                <div className="menu-dropdown" ref={optionsRef}>
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

        <div className="titlebar-center" data-tauri-drag-region={false}>
          <div className="auth-wrapper">
            {!currentUser ? (
              <>
                <button
                  className={`signin-btn auth-btn ${
                    showAuthMenu ? "active" : ""
                  }`}
                  onClick={() => {
                    setShowViewMenu(false);
                    setShowOptionsMenu(false);
                    setShowAuthMenu((prev) => !prev);
                  }}
                >
                  Sign In
                </button>

                {showAuthMenu && (
                  <div className="auth-dropdown" ref={authRef}>
                    <button
                      className="auth-item auth-register"
                      onClick={() => {
                        setShowAuthMenu(false);
                        setAuthPanelType("register");
                      }}
                    >
                      Register
                    </button>

                    <button
                      className="auth-item auth-login"
                      onClick={() => {
                        setShowAuthMenu(false);
                        setAuthPanelType("login");
                      }}
                    >
                      Log In
                    </button>
                  </div>
                )}
              </>
            ) : (
              <>
                <button
                  className={`user-btn ${showUserMenu ? "active" : ""}`}
                  onClick={() => {
                    setShowViewMenu(false);
                    setShowOptionsMenu(false);
                    setShowUserMenu((prev) => !prev);
                  }}
                >
                  <span className="user-bracket">&lt;</span>
                  <span className="user-name">{currentUser.username}</span>
                  <span className="user-bracket">&gt;</span>
                </button>

                {showUserMenu && (
                  <div className="user-dropdown" ref={userMenuRef}>
                    <button
                      className="user-menu-item logout"
                      onClick={async () => {
                        try {
                          await invoke("clear_tokens");
                          setCurrentUser(null);
                          setShowUserMenu(false);
                        } catch (error) {
                          console.error("Logout failed:", error);
                          const errorMsg =
                            typeof error === "string"
                              ? error
                              : "Failed to logout. Please try again.";
                          alert(errorMsg);
                        }
                      }}
                    >
                      Logout
                    </button>
                  </div>
                )}
              </>
            )}
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

      {authPanelType && (
        <div className="auth-panel">
          <div className="auth-panel-header">
            <button
              className="auth-panel-close"
              onClick={() => {
                setAuthPanelType(null);
              }}
            >
              ✕
            </button>
          </div>

          <div className="auth-panel-body">
            {authPanelType === "register" && (
              <RegisterForm
                onClose={() => setAuthPanelType(null)}
                onRegisterSuccess={(username: string, email: string) => {
                  setCurrentUser({ username, email });
                }}
              />
            )}
            {authPanelType === "login" && (
              <LoginForm
                onClose={() => setAuthPanelType(null)}
                onLoginSuccess={(username: string, email: string) => {
                  setCurrentUser({ username, email });
                }}
              />
            )}
          </div>
        </div>
      )}

      <main className="container"></main>
    </>
  );
}

export default App;
