export function validateEmailValue(email: string): string {
  if (!email) return "Email is required.";
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return "Invalid email format.";
  return "";
}

export function validateUsernameValue(username: string): string {
  if (!username) return "Username is required.";

  if (!/^@[a-z0-9]{3,16}$/.test(username)) {
    return "Username must contain 3–16 lowercase letters or digits.";
  }

  return "";
}


export function validatePasswordValue(password: string): string {
  if (!password) return "Password is required.";

  const strongRules =
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/;

  if (!strongRules.test(password)) {
    return "Password must be at least 8 characters and include lowercase, uppercase, number and special symbol.";
  }

  return "";
}


export function validateConfirmPasswordValue(password: string, confirm: string): string {
  if (password !== confirm) return "Passwords do not match.";
  return "";
}

export function getPasswordStrength(password: string): "weak" | "medium" | "strong" {
  if (password.length >= 16) return "strong";
  if (password.length >= 12) return "medium";
  return "weak";
}


