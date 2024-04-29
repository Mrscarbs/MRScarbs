use sqlx::mysql::MySqlPool;
use chrono::{Local, Duration};
use std::env;
use dotenv::dotenv;
use sqlx::Row;

pub struct TimeDiff {
    interval: String,
    api_id: i32,
}

#[async_trait::async_trait]
pub trait TimeDiffFuncs {
    async fn get_db_details(&self, pool: &MySqlPool) -> Option<i64>;
    async fn last_time(&self, pool: &MySqlPool) -> chrono::DateTime<Local>;
}

#[async_trait::async_trait]
impl TimeDiffFuncs for TimeDiff {
    async fn get_db_details(&self, pool: &MySqlPool) -> Option<i64> {
        let result = sqlx::query!(
            "CALL stp_get_limits(?, ?)",
            self.api_id,
            &self.interval
        )
        .fetch_one(pool)
        .await
        .ok()?;

        Some(result.try_get::<i64, _>(0).ok()?)
    }

    async fn last_time(&self, pool: &MySqlPool) -> chrono::DateTime<Local> {
        if let Some(minutes) = self.get_db_details(pool).await {
            Local::now() - Duration::minutes(minutes)
        } else {
            Local::now()
        }
    }
}

#[tokio::main]
async fn main() {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let pool = MySqlPool::connect(&database_url).await.expect("Failed to create pool.");

    let time_diff = TimeDiff {
        interval: "oneminute".to_string(),
        api_id: 1,
    };

    let last_time = time_diff.last_time(&pool).await;
    println!("Last time: {}", last_time);
}
