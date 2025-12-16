import { invoke } from "@tauri-apps/api/core";
import { useState } from "react";
import { useEffect } from "react";
import { useDebounce } from "../../hooks/useDebounce";
import {
  validateEmailValue,
  validateUsernameValue,
  validatePasswordValue,
  validateConfirmPasswordValue,
  getPasswordStrength,
} from "./validators";
const MAIL_MAX_LENGTH = 254;
const USERNAME_MAX_LENGTH = 16;
const PASSWORD_MAX_LENGTH = 72;

type RegisterFormProps = {
  onClose: () => void;
  onRegisterSuccess: (username: string, email: string) => void;
};

export default function RegisterForm({
  onClose,
  onRegisterSuccess,
}: RegisterFormProps) {
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("@");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [passwordStrength, setPasswordStrength] = useState<
    "weak" | "medium" | "strong"
  >("weak");
  const [showPassword, setShowPassword] = useState(false);
  const [shakeField, setShakeField] = useState<string | null>(null);
  const [tempEmailMessage, setTempEmailMessage] = useState("");
  const [tempUsernameMessage, setTempUsernameMessage] = useState("");
  const [tempPasswordMessage, setTempPasswordMessage] = useState("");
  const [tempConfirmMessage, setTempConfirmMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [errors, setErrors] = useState({
    email: "",
    username: "",
    password: "",
    confirmPassword: "",
    general: "",
  });
  const debouncedEmail = useDebounce(email, 500);
  const debouncedUsername = useDebounce(username, 500);
  function validateRegisterForm() {
    const emailError = validateEmailValue(email);
    const usernameError = validateUsernameValue(username);
    const passwordError = validatePasswordValue(password);
    const confirmPasswordError = validateConfirmPasswordValue(
      password,
      confirmPassword
    );

    setErrors({
      email: emailError,
      username: usernameError,
      password: passwordError,
      confirmPassword: confirmPasswordError,
      general: "",
    });

    return (
      !emailError && !usernameError && !passwordError && !confirmPasswordError
    );
  }

  async function handleSubmit() {
    if (!validateRegisterForm()) return;

    setIsSubmitting(true);

    try {
      const response = await invoke<{
        success: boolean;
        field_errors?: Array<{ field: string; message: string }>;
        message: string;
        access_token?: string;
        refresh_token?: string;
        user?: { id: string; username: string; email: string };
      }>("register_user", {
        payload: {
          email: email,
          username: username,
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
            onRegisterSuccess(response.user.username, response.user.email);
          }
          setIsSubmitting(false);
          setSuccessMessage("");
          onClose();
        }, 2000);
      } else {
        if (response.field_errors) {
          const newErrors: any = {};
          response.field_errors.forEach(({ field, message }) => {
            newErrors[field] = message;
          });
          setErrors((prev) => ({ ...prev, ...newErrors }));
        } else {
          setErrors((prev) => ({
            ...prev,
            general: response.message,
          }));
        }

        setIsSubmitting(false);
      }
    } catch (error) {
      console.error("Registration failed:", error);

      setErrors((prev) => ({
        ...prev,
        general:
          typeof error === "string"
            ? error
            : "An unexpected error occurred. Please try again.",
      }));

      setIsSubmitting(false);
    }
  }

  function showEmailTemp(msg: string) {
    setTempEmailMessage(msg);
    setTimeout(() => setTempEmailMessage(""), 2000);
  }

  function showUsernameTemp(msg: string) {
    setTempUsernameMessage(msg);
    setTimeout(() => setTempUsernameMessage(""), 2000);
  }

  function showPasswordTemp(msg: string) {
    setTempPasswordMessage(msg);
    setTimeout(() => setTempPasswordMessage(""), 2000);
  }

  function showConfirmTemp(msg: string) {
    setTempConfirmMessage(msg);
    setTimeout(() => setTempConfirmMessage(""), 1200);
  }

  function triggerShake(field: string) {
    setShakeField(field);

    setTimeout(() => {
      setShakeField(null);
    }, 250);
  }

  useEffect(() => {
    if (!debouncedEmail || validateEmailValue(debouncedEmail)) {
      return;
    }

    invoke<boolean | null>("check_email_available", { email: debouncedEmail })
      .then((available) => {
        if (available === false) {
          setErrors((prev) => ({
            ...prev,
            email: "This email is already taken.",
          }));
        }
      })
      .catch((err) => {
        console.error("Check email failed:", err);
        setErrors((prev) => ({
          ...prev,
          email: "Could not verify email. Please try again.",
        }));
      });
  }, [debouncedEmail]);

  useEffect(() => {
    if (!debouncedUsername || validateUsernameValue(debouncedUsername)) {
      return;
    }

    invoke<boolean | null>("check_username_available", {
      username: debouncedUsername,
    })
      .then((available) => {
        if (available === false) {
          setErrors((prev) => ({
            ...prev,
            username: "This username is already taken.",
          }));
        }
      })
      .catch((err) => {
        console.error("Check username failed:", err);
        setErrors((prev) => ({
          ...prev,
          username: "Could not verify username. Please try again.",
        }));
      });
  }, [debouncedUsername]);

  return (
    <div className="auth-form">
      <label className="input-label">
        Email
        <input
          disabled={isSubmitting}
          type="email"
          className={`input-field ${
            shakeField === "email" ? "input-shake input-error" : ""
          }`}
          value={email}
          onChange={(e) => {
            const value = e.target.value;

            if (value.length > MAIL_MAX_LENGTH) {
              showEmailTemp("Email cannot exceed 254 characters.");
              triggerShake("email");
              return;
            }

            setEmail(value);

            if (errors.email) {
              const msg = validateEmailValue(value);
              setErrors((prev) => ({ ...prev, email: msg }));
            }
          }}
          onBlur={() => {
            const normalized = email.trim().toLowerCase();
            setEmail(normalized);

            if (
              errors.email !== "This email is already taken." &&
              errors.email !== "Could not verify email. Please try again."
            ) {
              const msg = validateEmailValue(normalized);
              setErrors((prev) => ({ ...prev, email: msg }));
            }
          }}
          placeholder="you@example.com"
        />
      </label>

      {errors.email && !tempEmailMessage && (
        <p className="error-text">{errors.email}</p>
      )}

      {tempEmailMessage && <p className="temp-message">{tempEmailMessage}</p>}

      <label className="input-label">
        Username
        <input
          disabled={isSubmitting}
          type="text"
          className={`input-field ${
            shakeField === "username" ? "input-shake input-error" : ""
          }`}
          value={username}
          onChange={(e) => {
            let value = e.target.value.toLowerCase();

            if (!value.startsWith("@")) {
              value = "@" + value.replace(/@/g, "");
            }

            const usernameBody = value.slice(1);

            if (usernameBody.length > USERNAME_MAX_LENGTH) {
              showUsernameTemp(
                "Username can have at most 16 characters after @."
              );
              triggerShake("username");
              return;
            }

            const finalUsername = "@" + usernameBody;
            setUsername(finalUsername);
            if (errors.username) {
              const msg = validateUsernameValue(finalUsername);
              setErrors((prev) => ({ ...prev, username: msg }));
            }
          }}
          onBlur={() => {
            if (
              errors.username !== "This username is already taken." &&
              errors.username !== "Could not verify username. Please try again."
            ) {
              const msg = validateUsernameValue(username);
              setErrors((prev) => ({ ...prev, username: msg }));
            }
          }}
          placeholder="@username"
        />
      </label>

      {errors.username && !tempUsernameMessage && (
        <p className="error-text">{errors.username}</p>
      )}

      {tempUsernameMessage && (
        <p className="temp-message">{tempUsernameMessage}</p>
      )}

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

              const msg = validatePasswordValue(value);
              setErrors((prev) => ({ ...prev, password: msg }));

              if (!msg) {
                const strength = getPasswordStrength(value);
                setPasswordStrength(strength);
              } else {
                setPasswordStrength("weak");
              }

              // Revalidate confirm password if it has a value
              if (confirmPassword) {
                const confirmMsg = validateConfirmPasswordValue(
                  value,
                  confirmPassword
                );
                setErrors((prev) => ({ ...prev, confirmPassword: confirmMsg }));
              }
            }}
            onBlur={() => {
              const msg = validatePasswordValue(password);
              setErrors((prev) => ({ ...prev, password: msg }));
            }}
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

      {errors.password && !tempPasswordMessage && (
        <p className="error-text">{errors.password}</p>
      )}

      {tempPasswordMessage && (
        <p className="temp-message">{tempPasswordMessage}</p>
      )}

      {!errors.password && !tempPasswordMessage && password.length > 0 && (
        <p className={`password-strength ${passwordStrength}`}>
          {passwordStrength === "weak" && "Weak password"}
          {passwordStrength === "medium" && "Medium strength"}
          {passwordStrength === "strong" && "Strong password"}
        </p>
      )}

      <label className="input-label">
        Confirm Password
        <input
          disabled={isSubmitting}
          type="password"
          className={`input-field ${
            shakeField === "confirm" ? "input-shake input-error" : ""
          }`}
          value={confirmPassword}
          onChange={(e) => {
            const value = e.target.value;

            if (value.length > PASSWORD_MAX_LENGTH) {
              showConfirmTemp("Confirm password cannot exceed 72 characters.");
              triggerShake("confirm");
              return;
            }

            setConfirmPassword(value);

            const msg = validateConfirmPasswordValue(password, value);
            setErrors((prev) => ({ ...prev, confirmPassword: msg }));
          }}
          onBlur={() => {
            const msg = validateConfirmPasswordValue(password, confirmPassword);
            setErrors((prev) => ({ ...prev, confirmPassword: msg }));
          }}
          placeholder="********"
        />
      </label>

      {errors.confirmPassword && !tempConfirmMessage && (
        <p className="error-text">{errors.confirmPassword}</p>
      )}

      {tempConfirmMessage && (
        <p className="temp-message">{tempConfirmMessage}</p>
      )}

      <button
        className="auth-submit"
        onClick={handleSubmit}
        disabled={
          isSubmitting ||
          !!errors.email ||
          !!errors.username ||
          !!errors.password ||
          !!errors.confirmPassword ||
          !email ||
          !username ||
          !password ||
          !confirmPassword
        }
      >
        {isSubmitting ? <div className="spinner"></div> : "Create Account"}
      </button>
      {successMessage && <p className="success-text">{successMessage}</p>}
      {errors.general && (
        <p className="error-text-centered">{errors.general}</p>
      )}
    </div>
  );
}
