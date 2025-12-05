pub struct ApiConfig {
    pub base_url: String,
}

impl ApiConfig {
    pub fn new() -> Self {
        let base_url =
            std::env::var("API_BASE_URL").unwrap_or_else(|_| "http://localhost:3000".to_string());

        Self { base_url }
    }

    pub fn auth_url(&self, endpoint: &str) -> String {
        format!("{}/api/auth/{}", self.base_url, endpoint)
    }
}

impl Default for ApiConfig {
    fn default() -> Self {
        Self::new()
    }
}
