import { useState } from "react";
import {
  validateEmailValue,
  validateUsernameValue,
  validatePasswordValue,
  validateConfirmPasswordValue,
  getPasswordStrength,
} from "./validators";
const MAIL_MAX_LENGTH = 254;
const USERNAME_MAX_LENGTH = 16;
const PASSWORD_MAX_LENGTH = 128;

type RegisterFormProps = {
  onClose: () => void;
};

export default function RegisterForm({ onClose }: RegisterFormProps) {
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
  });

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
    });

    return (
      !emailError && !usernameError && !passwordError && !confirmPasswordError
    );
  }

  async function handleSubmit() {
    if (!validateRegisterForm()) return;

    setIsSubmitting(true);

    await new Promise((res) => setTimeout(res, 1200));

    setSuccessMessage("Account created successfully!");

    setTimeout(() => {
      setIsSubmitting(false);
      setSuccessMessage("");
      onClose();
    }, 1500);
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

            const msg = validateEmailValue(normalized);
            setErrors((prev) => ({ ...prev, email: msg }));
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

            const msg = validateUsernameValue(finalUsername);
            setErrors((prev) => ({ ...prev, username: msg }));
          }}
          onBlur={() => {
            const msg = validateUsernameValue(username);
            setErrors((prev) => ({ ...prev, username: msg }));
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
                showPasswordTemp("Password cannot exceed 128 characters.");
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
              showConfirmTemp("Confirm password cannot exceed 128 characters.");
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
    </div>
  );
}
