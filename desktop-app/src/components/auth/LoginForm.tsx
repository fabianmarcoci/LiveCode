import { invoke } from "@tauri-apps/api/core";
import { useState } from "react";

type LoginFormProps = {
  onClose: () => void;
};

export default function LoginForm({ onClose }: LoginFormProps) {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");

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
          className="input-field"
          value={identifier}
          onChange={(e) => {
            setIdentifier(e.target.value);
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
            className="input-field"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
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
            {showPassword ? "#" : "=A"}
          </span>
        </div>
      </label>

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
