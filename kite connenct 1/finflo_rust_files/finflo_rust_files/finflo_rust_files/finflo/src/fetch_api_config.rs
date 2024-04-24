use sqlx::{MySqlPool, Error};
use std::env;
use dotenv::dotenv;
use sqlx::Row;

#[derive(Debug)]
pub struct ApiConfig {
    pub api_key: Option<String>,
    pub secret_key: Option<String>,
    pub api_provider: Option<String>,
    pub access_token: Option<String>,
    pub last_purchase_time: Option<i64>,
    pub first_purchase_time: Option<i64>,
    pub historical: Option<i8>,
    pub instrument_type: Option<String>,
    pub api_id: i32,
}

pub async fn fetch_api_config_dynamic(api_id: i32) -> Result<ApiConfig, Error> {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let pool = MySqlPool::connect(&database_url).await?;

    let row = sqlx::query!("CALL stp_get_api_config(?)", api_id)
    .fetch_one(&pool)
    .await?;

    let api_config = ApiConfig {
        api_key: row.try_get(0)?,
        secret_key: row.try_get(1)?,
        api_provider: row.try_get(2)?,
        access_token: row.try_get(3)?,
        last_purchase_time: row.try_get(4)?,
        first_purchase_time: row.try_get(5)?,
        historical: row.try_get(6)?,
        instrument_type: row.try_get(7)?,
        api_id: row.try_get(8)?,
    };

    Ok(api_config)
}

// #[tokio::main]
// async fn main() {
//     match fetch_api_config_dynamic(1).await {
//         Ok(api_config) => {
//             println!("{:#?}", api_config);
//         },
//         Err(e) => println!("Failed to fetch API config: {}", e),
//     }
// }
