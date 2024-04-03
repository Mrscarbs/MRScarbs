use chrono::{NaiveDateTime, DateTime, Utc};
use polars::prelude::*;
use rusty_talib::{moving_average, exponential_moving_average};
use std::collections::HashMap;

struct Trade {
    entry_price: f64,
    exit_price: f64,
    entry_time: NaiveDateTime,
    exit_time: Option<NaiveDateTime>,
}

struct Backtester {
    data: Vec<(NaiveDateTime, f64, f64, f64, f64, f64)>,
    indicators: HashMap<String, HashMap<String, Vec<f64>>>,
    trades: Vec<Trade>,
}

impl Backtester {
    fn new(data: Vec<(NaiveDateTime, f64, f64, f64, f64, f64)>) -> Self {
        Backtester {
            data,
            indicators: HashMap::new(),
            trades: Vec::new(),
        }
    }

    fn calculate_indicator(&mut self, indicator_name: &str, price_type: &str, params: &[i32]) {
        let mut prices = Vec::new();
        for &(_, open, high, low, close, _) in &self.data {
            match price_type {
                "open" => prices.push(open),
                "high" => prices.push(high),
                "low" => prices.push(low),
                "close" => prices.push(close),
                _ => panic!("Unsupported price type: {}", price_type),
            }
        }

        let series = Series::new("", &prices);
        let indicator_values = match indicator_name {
            "sma" => moving_average(&series, Some(params[0] as u32)).unwrap(),
            "ema" => exponential_moving_average(&series, Some(params[0] as u32)).unwrap(),
            _ => panic!("Unsupported indicator: {}", indicator_name),
        };

        let indicator_values_vec: Vec<f64> = indicator_values.f64().unwrap().into_iter().flatten().collect();

        self.indicators
            .entry(indicator_name.to_string())
            .or_default()
            .insert(price_type.to_string(), indicator_values_vec);
    }

    fn backtest(&mut self, entry_condition: &dyn Fn(&Self, usize) -> bool, exit_condition: &dyn Fn(&Self, usize, &Trade) -> bool, take_profit: f64, stop_loss: f64) {
        for i in 0..self.data.len() {
            if entry_condition(self, i) {
                let entry_price = self.data[i].4; // Assuming close price for entry
                let trade = Trade {
                    entry_price,
                    exit_price: -1.0,
                    entry_time: self.data[i].0,
                    exit_time: None,
                };
                self.trades.push(trade);
            }

            if let Some(trade) = self.trades.last_mut() {
                if trade.exit_price == -1.0 {
                    let current_price = self.data[i].4; // Assuming close price for exit
                    if exit_condition(self, i, trade) || current_price >= trade.entry_price + take_profit || current_price <= trade.entry_price - stop_loss {
                        trade.exit_price = current_price;
                        trade.exit_time = Some(self.data[i].0);
                    }
                }
            }
        }
    }

    fn generate_statistics(&self) {
        println!("Total trades: {}", self.trades.len());
    }
}

fn main() {
    let data = vec![
        (DateTime::<Utc>::from_timestamp(1, 0).unwrap().naive_utc(), 100.0, 110.0, 90.0, 105.0, 1000.0),
    ];
    let mut backtester = Backtester::new(data);
    backtester.calculate_indicator("sma", "close", &[20]);
    backtester.calculate_indicator("ema", "open", &[12]);

    let entry_condition = |backtester: &Backtester, i: usize| -> bool {
        let close = backtester.data[i].4;
        let sma = backtester.indicators.get("sma").unwrap().get("close").unwrap()[i];
        close > sma
    };

    let exit_condition = |_backtester: &Backtester, _i: usize, _trade: &Trade| -> bool {
        false
    };

    let take_profit = 10.0;
    let stop_loss = 5.0;

    backtester.backtest(&entry_condition, &exit_condition, take_profit, stop_loss);

    backtester.generate_statistics();
}
