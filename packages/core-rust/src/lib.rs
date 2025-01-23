use serde_json::{Result, Value};
use std::fs::File;
use std::io::prelude::*;

pub fn greet(name: &str) -> String {
    format!("Hello, {}!", name)
}

pub fn read_config() -> Result<()> {
    let data = r#"
        {
            "name":"hello"
        }
    "#;

    let v: Value = serde_json::from_str(data)?;
    println!("{}", v["name"]);

    Ok(())
}

pub fn test() -> std::io::Result<()> {}
