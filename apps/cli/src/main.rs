use clap::Parser;
use std::fs::File;
use std::io::prelude::*;
// 重构 ordo cli 工具，使用 rust 减少打包体积，优化性能

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    /// Name of the person to greet
    #[arg(short, long)]
    name: String,

    /// Number of times to greet
    #[arg(short, long, default_value_t = 1)]
    count: u8,
}

fn main() -> std::io::Result<()> {

}
