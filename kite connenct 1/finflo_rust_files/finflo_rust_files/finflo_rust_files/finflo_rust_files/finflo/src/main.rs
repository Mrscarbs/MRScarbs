use std::env;
use std::time::{SystemTime, UNIX_EPOCH};
use dotenv::dotenv;
use sqlx::{mysql::MySqlPool, Row};
use chrono::Local;


pub mod quant_indicators;
pub mod fetch_api_config;
pub mod enums;

async fn quant_insert(instrument_token_list: Vec<i32>, time_frame:&str, interval:&str){
    
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let pool = MySqlPool::connect(&database_url).await.expect("Failed to create pool.");

    for (_index, element) in instrument_token_list.iter().enumerate() { 

        let current_time = SystemTime::now();
        let epoch_timestamp = current_time.duration_since(UNIX_EPOCH)
                                            .expect("Time went backwards")
                                            .as_secs();

        let indicators = quant_indicators::QuantIndicators {
             ticker_id: *element,
             time_frame: time_frame.to_string(),
         };

        let local_datetime = Local::now();
        let naive_datetime = local_datetime.naive_utc();
        let to_date = naive_datetime.format("%Y-%m-%d %H:%M:%S").to_string();

        let sharpe = indicators.sharpe(&to_date, interval.to_string(), quant_indicators::PriceType::Close, 0.03)
                                                        .await
                                                        .expect("Failed to calculate Sharpe ratio");
        let sortino = indicators.sortino(&to_date, interval.to_string(), quant_indicators::PriceType::Close, 0.03, 0.00)
                                                            .await
                                                            .expect("Failed to calculate Sortino ratio");
         

        let _result = sqlx::query!(
                "CALL stp_insert_or_Update_QuantStats(?, ?, ?, ?)",
                sortino,
                element,
                sharpe,
                epoch_timestamp
            )
            .fetch_one(&pool)
            .await
            .ok();
    }
}

#[tokio::main]
async fn main() {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let pool = MySqlPool::connect(&database_url).await.expect("Failed to create pool.");

    let instrument_tokens = sqlx::query("CALL stp_get_tbl_current_ltp()")
                                    .fetch_all(&pool)
                                    .await
                                    .expect("Failed to fetch instrument tokens");

    let instrument_token_list: Vec<i32> = instrument_tokens.iter().map(|row| row.try_get(0)
                                                            .expect("error while retrieving tokens list"))
                                                            .collect();

    println!("{:?}",instrument_token_list);

    let time_frame = String::from("minute");
    let interval = String::from("oneminute");

    quant_insert(instrument_token_list, &time_frame, &interval).await;
}
