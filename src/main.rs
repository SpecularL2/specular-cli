use clap::Parser;
use envconfig::Envconfig;

use crate::workspace::WorkspaceManager;
mod workspace;

#[derive(Parser)]
#[command(author, version, about, long_about = None, arg_required_else_help = true)]
struct Args {
    /// Scope of the operation
    scope: String,
    /// Command options
    command: String,
    /// Object name
    name: String,
}

#[derive(Envconfig)]
pub struct Config {
    #[envconfig(from = "SPC_PATH", default = "~/.spc")]
    pub spc_path: String,
}

struct App {
    workspace_manager: WorkspaceManager,
}

impl App {
    pub fn new(cfg: Config) -> App {
        let workspace_manager = WorkspaceManager::new(cfg);
        App {
            workspace_manager: workspace_manager,
        }
    }

    fn action(&self) -> Result<(), Box<dyn std::error::Error>> {
        let args = Args::parse();

        println!("Scope: {}", args.scope);
        println!("Command: {}", args.command);
        println!("Name: {}", args.name);

        match args.scope.as_str() {
            "workspace" => self.workspace_manager.manage(args.command, args.name),
            _ => Err("unknown scope")?,
        }
    }
}

fn main() {
    let cfg = Config::init_from_env().unwrap();
    let app: App = App::new(cfg);

    match app.action() {
        Ok(_) => println!("All good!"),
        Err(err) => eprintln!("Error: {}", err),
    }
}
