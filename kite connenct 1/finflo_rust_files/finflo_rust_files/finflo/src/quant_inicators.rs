use sqlx::mysql::MySqlPool;
use chrono::{Local, Duration};
use std::env;
use dotenv::dotenv;
use sqlx::Row;
extern crate kiteconnect;
extern crate serde_json as json;
use kiteconnect::connect::KiteConnect;
extern crate reqwest;
extern crate serde_json as json;
use reqwest::header::{HeaderMap, HeaderValue, AUTHORIZATION};
use reqwest::Error;
use serde_json::Value;

pub struct quant_indicators{
    ticker : String,
    time_frame : String,
}

impl quant_indicators{
    
    async fn get_historical_data(api_key: &str, access_token: &str, instrument_token: u32, from_date: &str, to_date: &str, interval: &str) -> Result<Value, Error> {
        let url = format!("https://api.kite.trade/instruments/historical/{}/{}?from={}&to={}", instrument_token, interval, from_date, to_date);
        let mut headers = HeaderMap::new();
        headers.insert("X-Kite-Version", HeaderValue::from_static("3"));
        headers.insert(AUTHORIZATION, HeaderValue::from_str(&format!("token {}:{}", api_key, access_token))?);
    
        let client = reqwest::Client::new();
        let response = client.get(&url)
            .headers(headers)
            .send()
            .await?;
    
        let historical_data = response.json::<Value>().await?;
        println!(historical_data)
    }
    
    fn sharpe(self)->f64{
        
    }
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    let api_key = "your_api_key";
    let access_token = "your_access_token";
    let instrument_token = 5633; // Example instrument token
    let from_date = "2017-12-15+09:15:00";
    let to_date = "2017-12-15+09:20:00";
    let interval = "minute";

    let historical_data = get_historical_data(api_key, access_token, instrument_token, from_date, to_date, interval).await?;
    while true{
    println!("{:#?}", historical_data);
    }
    Ok(())
}
