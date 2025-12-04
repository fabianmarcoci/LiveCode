use super::models::{RegisterRequest, RegisterResponse};
use sqlx::PgPool;
use super::super::password::hash_password;

pub async fn check_field_available_internal(
    field_name: &str,
    value: String,
    pool: &PgPool,
) -> Option<bool> {
    let query = format!(
        "SELECT id FROM users WHERE {} = $1 LIMIT 1",
        field_name
    );

    let exists = sqlx::query(&query)
        .bind(&value)
        .fetch_optional(pool)
        .await
        .ok()?
        .is_some();

    Some(!exists)
}

pub async fn register_user_internal(
    payload: RegisterRequest,
     pool: &PgPool,
) -> Result<RegisterResponse, String> {
    let mut field_errors: Vec<(String, String)> = vec![];

    let email_exists = sqlx::query!(
        "SELECT id FROM users WHERE email = $1 LIMIT 1",
        payload.email
    )
    .fetch_optional(pool)
    .await
    .map_err(|_| "Database error during email check.")?
    .is_some();

    if email_exists {
        field_errors.push(("email".into(), "This email is already taken.".into()));
    }

    let username_exists = sqlx::query!(
        "SELECT id FROM users WHERE username = $1 LIMIT 1",
        payload.username
    )
    .fetch_optional(pool)
    .await
    .map_err(|_| "Database error during username check.")?
    .is_some();

    if username_exists {
        field_errors.push(("username".into(), "This username is already taken.".into()));
    }

    if !field_errors.is_empty() {
        return Ok(RegisterResponse {
            success: false,
            field_errors: Some(field_errors),
            message: "Account could not be created.".into(),
        });
    }

    let password_hash = hash_password(&payload.password)?;

    sqlx::query!(
        "INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)",
        payload.email,
        payload.username,
        password_hash
    )
    .execute(pool)
    .await
    .map_err(|_| "An unexpected error occurred. Please try again.")?;

    Ok(RegisterResponse {
        success: true,
        field_errors: None,
        message: "Your account has been created successfully.".into(),
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_register_user_success() {
        dotenvy::dotenv().ok();
        let pool = sqlx::postgres::PgPoolOptions::new()
            .connect(&std::env::var("DATABASE_URL").unwrap())
            .await
            .unwrap();

        let payload = RegisterRequest {
            email: "test@example1.com".to_string(),
            username: "@testuser1".to_string(),
            password: "TestPass123!".to_string(),
        };

        let result = register_user_internal(payload, &pool).await;

        assert!(result.is_ok());
        let response = result.unwrap();
        assert_eq!(response.success, true);
    }
}