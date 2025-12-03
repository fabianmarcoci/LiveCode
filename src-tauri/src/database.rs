use sqlx::{Pool, Postgres};
use std::env;

pub async fn init_db() -> Pool<Postgres> {
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set in .env file");

    let pool = sqlx::postgres::PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .unwrap_or_else(|e| {
            eprintln!("Failed to connect to PostgreSQL: {}", e);
            std::process::exit(1);
        });

    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .expect("Failed to run migrations");

    pool
}
