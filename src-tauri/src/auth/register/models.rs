use serde::{Deserialize, Serialize};


#[derive(Deserialize)]
pub struct RegisterRequest {
    pub email: String,
    pub username: String,
    pub password: String,
}

#[derive(Serialize)]
pub struct RegisterResponse {
    pub success: bool,
    pub field_errors: Option<Vec<(String, String)>>, 
    pub message: String,
}
