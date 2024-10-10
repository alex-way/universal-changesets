use clap::Parser;
use cli::{
    add_changeset, get_version, version_command, AddCommand, Cli, Commands, PreviewCommands,
};

mod changelog;
mod changeset;
mod cli;
mod config;
mod plugin;
mod utils;

fn main() {
    let cli = Cli::parse();

    match &cli.command {
        Some(Commands::Add(command)) => add_changeset(command),
        Some(Commands::Version(_command)) => version_command(),
        Some(Commands::Get(_)) => {
            get_version();
        }
        Some(Commands::Preview(command)) => match &command.preview_commands {
            PreviewCommands::Version(_) => {
                get_version();
            }
        },
        None => add_changeset(&AddCommand {
            increment_type: None,
            message: None,
        }),
    }
}
