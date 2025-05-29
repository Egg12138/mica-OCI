mod commands;
use clap::{crate_version, CommandFactory, Parser};

#[macro_export]
macro_rules! rsmica_version {
    () => {
       concat!(
        "version ",
        crate_version!(),
       )
    };
}

#[derive(Parser, Debug)]
enum SubCommand {
    // Minimal commands set required by OCI-spec
    #[clap(flatten)]
    Minimals(Box<commands::MinimalCmd>),
    Common(Box<commands::CommonCmd>),

    // Mica extensions
}

fn main() {
    println!("Hello, world!");
}
