use chrono::NaiveDateTime;
use serde::Deserialize;
use csv::ReaderBuilder;
use std::collections::HashMap;
use std::cell::RefCell;
use std::rc::Rc;
use std::fmt::{self, Display, Formatter};
use yata::methods::{EMA, SMA};
use yata::prelude::*;

#[derive(Debug, Deserialize)]
struct CsvRecord {
    Date: String,
    Open: f64,
    High: f64,
    Low: f64,
    Close: f64,
    Volume: f64,
}

struct MarketData {
    dates: Vec<NaiveDateTime>,
    opens: Vec<f64>,
    highs: Vec<f64>,
    lows: Vec<f64>,
    closes: Vec<f64>,
    volumes: Vec<f64>,
}

impl MarketData {
    fn new() -> Self {
        MarketData {
            dates: Vec::new(),
            opens: Vec::new(),
            highs: Vec::new(),
            lows: Vec::new(),
            closes: Vec::new(),
            volumes: Vec::new(),
        }
    }

    fn push(&mut self, date: NaiveDateTime, open: f64, high: f64, low: f64, close: f64, volume: f64) {
        self.dates.push(date);
        self.opens.push(open);
        self.highs.push(high);
        self.lows.push(low);
        self.closes.push(close);
        self.volumes.push(volume);
    }
}

fn read_csv_data(path: &str) -> Result<MarketData, csv::Error> {
    let mut reader = ReaderBuilder::new().from_path(path)?;
    let mut data = MarketData::new();

    for result in reader.deserialize() {
        let record: CsvRecord = result?;
        let date = NaiveDateTime::parse_from_str(&record.Date, "%d-%m-%Y %H:%M")
            .expect("Failed to parse date");
        data.push(date, record.Open, record.High, record.Low, record.Close, record.Volume);
    }

    Ok(data)
}

#[derive(Debug, PartialEq)]
enum PriceType {
    Open,
    High,
    Low,
    Close,
}

impl Display for PriceType {
    fn fmt(&self, f: &mut Formatter) -> fmt::Result {
        match self {
            PriceType::Open => write!(f, "Open"),
            PriceType::High => write!(f, "High"),
            PriceType::Low => write!(f, "Low"),
            PriceType::Close => write!(f, "Close"),
        }
    }
}

struct Trade {
    entry_price: f64,
    exit_price: f64,
    entry_time: NaiveDateTime,
    exit_time: Option<NaiveDateTime>,
    profit_loss: f64,
    margin_used: f64,
}

struct Backtester {
    data: MarketData,
    indicators: Rc<RefCell<HashMap<String, HashMap<String, Vec<f64>>>>>,
    trades: Rc<RefCell<Vec<Trade>>>,
    initial_margin: f64,
    commission_rate: f64,
    equity: Vec<f64>,
    margin_percent_per_trade: f64,
    leverage: f64,
}

impl Backtester {
    fn new(data: MarketData, initial_margin: f64, commission_rate: f64, margin_percent_per_trade: f64, leverage: f64) -> Self {
        Backtester {
            data,
            indicators: Rc::new(RefCell::new(HashMap::new())),
            trades: Rc::new(RefCell::new(Vec::new())),
            initial_margin,
            commission_rate,
            equity: vec![initial_margin],
            margin_percent_per_trade,
            leverage,
        }
    }

    fn calculate_indicator(&self, indicator_name: &str, price_type: PriceType) {
        let prices = match price_type {
            PriceType::Open => &self.data.opens,
            PriceType::High => &self.data.highs,
            PriceType::Low => &self.data.lows,
            PriceType::Close => &self.data.closes,
        };

        let indicator_values = match indicator_name {
            "sma" => {
                let period = self.indicators.borrow().get(indicator_name).and_then(|params| params.get("period")).and_then(|v| v.first()).copied().unwrap_or(20.0) as usize;
                let mut sma = SMA::new(period.try_into().unwrap(), &prices[0]).unwrap();
                prices.iter().map(|&price| sma.next(&price)).collect::<Vec<f64>>()
            },
            "ema" => {
                let period = self.indicators.borrow().get(indicator_name).and_then(|params| params.get("period")).and_then(|v| v.first()).copied().unwrap_or(12.0) as usize;
                let mut ema = EMA::new(period.try_into().unwrap(), &prices[0]).unwrap();
                prices.iter().map(|&price| ema.next(&price)).collect::<Vec<f64>>()
            }
            _ => panic!("Unsupported indicator: {}", indicator_name),
        };

        self.indicators.borrow_mut().entry(indicator_name.to_string()).or_default().insert(price_type.to_string(), indicator_values);
    }

    fn backtest(
        &mut self,
        entry_conditions: &[Box<dyn Fn(&Backtester, usize) -> bool>],
        exit_conditions: &[Box<dyn Fn(&Backtester, &Trade, usize) -> bool>],
        take_profit: f64,
        stop_loss: f64,
    ) {
        let mut current_equity = self.initial_margin;
        let mut trades = self.trades.borrow_mut();
        let indicators = self.indicators.borrow();

        for i in 0..self.data.dates.len() {
            let margin_per_trade = current_equity * self.margin_percent_per_trade / 100.0;
            let entry_price = self.data.closes[i]; // Assuming close price for entry
            let position_size = margin_per_trade * self.leverage / entry_price;

            if current_equity > 0.0 && entry_conditions.iter().all(|condition| condition(self, i)) {
                let entry_commission = entry_price * position_size * self.commission_rate / 100.0;

                if current_equity > entry_commission + margin_per_trade {
                    let trade = Trade {
                        entry_price,
                        exit_price: 0.0,
                        entry_time: self.data.dates[i],
                        exit_time: None,
                        profit_loss: -entry_commission,
                        margin_used: margin_per_trade,
                    };
                    trades.push(trade);
                    current_equity -= entry_commission + margin_per_trade;
                }
            }

            for trade in trades.iter_mut().filter(|t| t.exit_price == 0.0) {
                let current_price = self.data.closes[i]; // Assuming close price for exit
                let current_time = self.data.dates[i];

                if current_time > trade.entry_time {
                    if exit_conditions.iter().any(|condition| condition(self, trade, i))
                        || current_price >= trade.entry_price + take_profit
                        || current_price <= trade.entry_price - stop_loss
                    {
                        let exit_commission = current_price * position_size * self.commission_rate / 100.0;
                        trade.exit_price = current_price;
                        trade.exit_time = Some(current_time);
                        let profit = (trade.exit_price - trade.entry_price) * position_size;
                        trade.profit_loss += profit - exit_commission;
                        current_equity += trade.margin_used + profit - exit_commission;
                    }
                }
            }
        }

        self.equity.push(current_equity);
    }

    fn generate_statistics(&self) {
        let trades = self.trades.borrow();
        println!("Total trades: {}", trades.len());
        println!("Final equity: {}", self.equity.last().unwrap());
    }
}

fn main() {
    let csv_path = "C:/finflo/MRScarbs/kite connenct 1/finflo_rust_files/finflo_rust_files/backtester/backtester/btc.csv";
    let data = read_csv_data(csv_path).expect("Failed to read CSV data");

    let mut backtester = Backtester::new(data, 10000.0, 0.04, 10.0, 10.0);  // 10% of equity per trade, 10x leverage

    backtester.indicators.borrow_mut().insert(
        "sma".to_string(),
        HashMap::from([("period".to_string(), vec![20.0])]),
    );
    backtester.indicators.borrow_mut().insert(
        "ema".to_string(),
        HashMap::from([("period".to_string(), vec![12.0])]),
    );

    backtester.calculate_indicator("sma", PriceType::Close);
    backtester.calculate_indicator("ema", PriceType::Close);

    let entry_conditions: Vec<Box<dyn Fn(&Backtester, usize) -> bool>> = vec![
        Box::new(|backtester, i| {
            let indicators = backtester.indicators.borrow();
            let sma = indicators.get("sma").and_then(|m| m.get("Close")).and_then(|v| v.get(i)).unwrap_or(&0.0);
            let close = backtester.data.closes[i];
            close > *sma
        }),
    ];

    let exit_conditions: Vec<Box<dyn Fn(&Backtester, &Trade, usize) -> bool>> = vec![
        Box::new(|backtester, trade, i| {
            let indicators = backtester.indicators.borrow();
            let sma = indicators.get("sma").and_then(|m| m.get("Close")).and_then(|v| v.get(i)).unwrap_or(&0.0);
            let ema = indicators.get("ema").and_then(|m| m.get("Close")).and_then(|v| v.get(i)).unwrap_or(&0.0);
            *sma > *ema
        }),
    ];

    backtester.backtest(&entry_conditions, &exit_conditions, 10.0, 5.0);

    backtester.generate_statistics();
}
