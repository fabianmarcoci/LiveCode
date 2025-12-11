import { invoke } from "@tauri-apps/api/core";
import { useState } from "react";

const IDENTIFIER_MAX_LENGTH = 255;
const PASSWORD_MAX_LENGTH = 72;

type LoginFormProps = {
  onClose: () => void;
  onLoginSuccess: (username: string, email: string) => void;
};

export default function LoginForm({ onClose, onLoginSuccess }: LoginFormProps) {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [shakeField, setShakeField] = useState<string | null>(null);
  const [tempIdentifierMessage, setTempIdentifierMessage] = useState("");
  const [tempPasswordMessage, setTempPasswordMessage] = useState("");

  function triggerShake(field: string) {
    setShakeField(field);
    setTimeout(() => setShakeField(null), 250);
  }

  function showIdentifierTemp(msg: string) {
    setTempIdentifierMessage(msg);
    setTimeout(() => setTempIdentifierMessage(""), 2000);
  }

  function showPasswordTemp(msg: string) {
    setTempPasswordMessage(msg);
    setTimeout(() => setTempPasswordMessage(""), 2000);
  }

  async function handleSubmit() {
    if (!identifier || !password) {
      setErrorMessage("Please fill in all fields.");
      return;
    }

    setIsSubmitting(true);
    setErrorMessage("");

    try {
      const response = await invoke<{
        success: boolean;
        message: string;
        access_token?: string;
        refresh_token?: string;
        user?: {
          id: string;
          username: string;
          email: string;
        };
      }>("login_user", {
        payload: {
          identifier: identifier,
          password: password,
        },
      });

      if (response.success) {
        if (response.access_token && response.refresh_token) {
          await invoke("save_tokens", {
            accessToken: response.access_token,
            refreshToken: response.refresh_token,
          });
        }

        setSuccessMessage(response.message);
        setTimeout(() => {
          if (response.user) {
            onLoginSuccess(response.user.username, response.user.email);
          }
          setIsSubmitting(false);
          setSuccessMessage("");
          onClose();
        }, 2000);
      } else {
        setErrorMessage(response.message);
        setIsSubmitting(false);
      }
    } catch (error) {
      console.error("Login failed:", error);

      setErrorMessage(
        typeof error === "string"
          ? error
          : "An unexpected error occurred. Please try again."
      );

      setIsSubmitting(false);
    }
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === "Enter" && !isSubmitting && identifier && password) {
      handleSubmit();
    }
  }

  return (
    <div className="auth-form">
      <label className="input-label">
        Email / Username
        <input
          disabled={isSubmitting}
          type="text"
          className={`input-field ${
            shakeField === "identifier" ? "input-shake input-error" : ""
          }`}
          value={identifier}
          onChange={(e) => {
            const value = e.target.value;

            if (value.length > IDENTIFIER_MAX_LENGTH) {
              showIdentifierTemp("Email or username cannot exceed 255 characters.");
              triggerShake("identifier");
              return;
            }

            setIdentifier(value);
            if (errorMessage) {
              setErrorMessage("");
            }
          }}
          onKeyDown={handleKeyDown}
          placeholder="you@example.com or @username"
        />
      </label>

      <label className="input-label">
        Password
        <div className="input-with-icon">
          <input
            disabled={isSubmitting}
            type={showPassword ? "text" : "password"}
            className={`input-field ${
              shakeField === "password" ? "input-shake input-error" : ""
            }`}
            value={password}
            onChange={(e) => {
              const value = e.target.value;

              if (value.length > PASSWORD_MAX_LENGTH) {
                showPasswordTemp("Password cannot exceed 72 characters.");
                triggerShake("password");
                return;
              }

              setPassword(value);
              if (errorMessage) {
                setErrorMessage("");
              }
            }}
            onKeyDown={handleKeyDown}
            placeholder="********"
          />

          <span
            className="password-toggle"
            onMouseDown={(e) => e.preventDefault()}
            onClick={() => setShowPassword((prev) => !prev)}
          >
            {showPassword ? "⌣" : "👁"}
          </span>
        </div>
      </label>

      {tempIdentifierMessage && <p className="temp-message">{tempIdentifierMessage}</p>}
      {tempPasswordMessage && <p className="temp-message">{tempPasswordMessage}</p>}

      <button
        className="auth-submit"
        onClick={handleSubmit}
        disabled={isSubmitting || !identifier || !password}
      >
        {isSubmitting ? <div className="spinner"></div> : "Sign In"}
      </button>
      {successMessage && <p className="success-text">{successMessage}</p>}
      {errorMessage && <p className="error-text-centered">{errorMessage}</p>}
    </div>
  );
}
