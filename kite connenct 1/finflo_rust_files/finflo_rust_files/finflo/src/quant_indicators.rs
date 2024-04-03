use anyhow::Error;

use dotenv::dotenv;
use reqwest::header::{HeaderMap, HeaderValue, AUTHORIZATION};
use serde_json::Value;
use sqlx::mysql::MySqlPool;
use std::env;

use crate::fetch_api_config;
use crate::enums;

#[derive(Clone)] // Implement Clone trait for PriceType
pub enum PriceType {
    Open,
    High,
    Low,
    Close,
}

pub struct QuantIndicators {
    pub ticker_id: i32,
    pub time_frame: String,
}

impl QuantIndicators {
    pub async fn get_historical_data(&self, to_date: &str, interval: String) -> Result<Value, Error> {
        dotenv().ok();
        
        let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
        let pool = MySqlPool::connect(&database_url).await.expect("Failed to create pool.");

        let api_config_result = fetch_api_config::fetch_api_config_dynamic(1).await;
        let api_config = api_config_result?;

        let api_key: String = api_config.api_key.unwrap();
        let access_token: String = api_config.access_token.unwrap();

        let time_diff = enums::TimeDiff {
            interval: interval.to_string(),
            api_id: 1,
        };

        let from_date_str = time_diff.last_time(&pool).await;
        
        let from_date = from_date_str.format("%Y-%m-%d %H:%M:%S").to_string();

        let url = format!("https://api.kite.trade/instruments/historical/{}/{}?from={}&to={}", self.ticker_id, self.time_frame, from_date, to_date);
        
        let mut headers = HeaderMap::new();
        headers.insert("X-Kite-Version", HeaderValue::from_static("3"));
        headers.insert(AUTHORIZATION, HeaderValue::from_str(&format!("token {}:{}", api_key, access_token)).expect("invalid header value"));

        let client = reqwest::Client::new();
        let response = client.get(&url)
            .headers(headers)
            .send()
            .await?;

        let response_text = response.text().await?;
        
        let response_json: Value = serde_json::from_str(&response_text)?;

        Ok(response_json)
    }

    pub async fn sharpe(&self, to_date: &str, interval: String, price_type: PriceType, risk_free_rate: f64) -> Result<f64, Error> {
        let historical_data = self.get_historical_data(to_date, interval.clone()).await?; // Clone interval here
        
        let mut prices = Vec::new();

        if let Some(data) = historical_data["data"]["candles"].as_array() {
            for entry in data {
                if let Some(entry_array) = entry.as_array() {
                    if entry_array.len() >= 5 {
                        match price_type {
                            PriceType::Open => prices.push(entry_array[1].as_f64().unwrap_or_default()),
                            PriceType::High => prices.push(entry_array[2].as_f64().unwrap_or_default()),
                            PriceType::Low => prices.push(entry_array[3].as_f64().unwrap_or_default()),
                            PriceType::Close => prices.push(entry_array[4].as_f64().unwrap_or_default()),
                        }
                    }
                }
            }
        }

        if prices.is_empty() {
            return Err(Error::msg("No price data found."));
        }

        let returns: Vec<f64> = prices.windows(2).map(|w| w[1] - w[0]).collect();

        let avg_return = returns.iter().sum::<f64>() / returns.len() as f64;
        let std_dev = (returns.iter().map(|&r| (r - avg_return).powi(2)).sum::<f64>() / (returns.len() - 1) as f64).sqrt();

        let sharpe_ratio = (avg_return - risk_free_rate) / std_dev;

        Ok(sharpe_ratio)
    }

    pub async fn sortino(&self, to_date: &str, interval: String, price_type: PriceType, risk_free_rate: f64, target_return: f64) -> Result<f64, anyhow::Error> {
        let historical_data = self.get_historical_data(to_date, interval).await?; // No need to clone interval here as it's the last use
        
        let mut prices = Vec::new();

        if let Some(data) = historical_data["data"]["candles"].as_array() {
            for entry in data {
                if let Some(entry_array) = entry.as_array() {
                    if entry_array.len() >= 5 {
                        match price_type {
                            PriceType::Open => prices.push(entry_array[1].as_f64().unwrap_or_default()),
                            PriceType::High => prices.push(entry_array[2].as_f64().unwrap_or_default()),
                            PriceType::Low => prices.push(entry_array[3].as_f64().unwrap_or_default()),
                            PriceType::Close => prices.push(entry_array[4].as_f64().unwrap_or_default()),
                        }
                    }
                }
            }
        }

        if prices.is_empty() {
            return Err(anyhow::Error::msg("No price data found."));
        }

        // Calculate returns
        let returns: Vec<f64> = prices.windows(2).map(|w| w[1] - w[0]).collect();

        // Calculate excess returns over the target return
        let excess_returns: Vec<f64> = returns.iter().map(|&r| r - target_return).collect();

        // Filter for negative excess returns to focus on downside risk
        let negative_excess_returns: Vec<f64> = excess_returns.iter().filter(|&&r| r < 0.0).cloned().collect();

        // Calculate the average excess return
        let avg_excess_return = excess_returns.iter().sum::<f64>() / returns.len() as f64;

        // Calculate the downside deviation (standard deviation of negative excess returns)
        let downside_deviation = if !negative_excess_returns.is_empty() {
            (negative_excess_returns.iter().map(|&r| r.powi(2)).sum::<f64>() / negative_excess_returns.len() as f64).sqrt()
        } else {
            0.0
        };

        // Calculate the Sortino ratio
        let sortino_ratio = if downside_deviation != 0.0 {
            (avg_excess_return - risk_free_rate) / downside_deviation
        } else {
            0.0
        };

        Ok(sortino_ratio)
    }
    
}

// #[tokio::main]
// async fn main() -> Result<(), Error> {
//     dotenv().ok();

//     let instrument_token = 5633; // Example instrument token
//     let local_datetime = Local::now();
//     let naive_datetime = local_datetime.naive_utc();
//     let to_date = naive_datetime.format("%Y-%m-%d %H:%M:%S").to_string();
//     let risk_free_rate = 0.03; // Example risk-free rate
//     let target_return = 0.03;
    
//     let indicators = QuantIndicators {
//         ticker_id: instrument_token,
//         time_frame: String::from("minute"),
//     };

//     let price_type = PriceType::Close; // Example price type
//     let interval = String::from("oneminute");
//     let sharpe_ratio = indicators.sharpe(&to_date, interval.clone(), price_type.clone(), risk_free_rate).await?; // Clone interval and price_type here
//     // Prefix unused variable with an underscore or remove it if not used
//     let _sortino_ratio = indicators.sortino(&to_date, interval, price_type, risk_free_rate, target_return).await?;

//     println!("Calculated Sharpe Ratio: {}", sharpe_ratio);

//     Ok(())
// }
