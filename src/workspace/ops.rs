use kdam::term::Colorizer;
use std::error::Error;
use std::fs;
use std::io;
use std::path::PathBuf;
use std::process;
use walkdir::WalkDir;

use crate::Config;

use super::output;
use super::parser;
use super::requests;

static DEFAULT_WORKSPACE_URL: &str =
    "https://github.com/SpecularL2/specular/tree/develop/config/local_devnet";

pub struct WorkspaceManager {
    default_workspace_url: String,
    workspaces_path: PathBuf,
    default_workspace_path: PathBuf,
}

impl WorkspaceManager {
    pub fn new(cfg: Config) -> WorkspaceManager {
        let spc_path: String = cfg.spc_path.clone();
        let home_dir: PathBuf = home::home_dir().unwrap();

        let abs_path = spc_path.replace("~", &home_dir.to_str().unwrap());
        let abs_base_path = PathBuf::from(abs_path);

        let workspaces_path: PathBuf = PathBuf::from(abs_base_path).join("workspaces");
        let default_workspace_path: PathBuf = workspaces_path.clone().join("default");
        WorkspaceManager {
            workspaces_path: workspaces_path,
            default_workspace_path: default_workspace_path,
            default_workspace_url: DEFAULT_WORKSPACE_URL.to_string(),
        }
    }

    fn cp(&self, in_dir: PathBuf, out_dir: PathBuf) -> Result<(), Box<dyn Error>> {
        if !in_dir.exists() {
            return Err("First, please run: spc workspace download default".into());
        }

        for entry in WalkDir::new(&in_dir) {
            let entry = entry?;

            let from = entry.path();
            let to = out_dir.join(from.strip_prefix(&in_dir)?);
            println!("\tcopy {} => {}", from.display(), to.display());

            // create directories
            if entry.file_type().is_dir() {
                if let Err(e) = fs::create_dir(to) {
                    match e.kind() {
                        io::ErrorKind::AlreadyExists => {}
                        _ => return Err(e.into()),
                    }
                }
            }
            // copy files
            else if entry.file_type().is_file() {
                fs::copy(from, to)?;
            }
            // ignore the rest
            else {
                eprintln!("copy: ignored symlink {}", from.display());
            }
        }
        Ok(())
    }

    fn rm(&self, rm_dir: PathBuf) -> Result<(), Box<dyn Error>> {
        print!("removing: {:?}", rm_dir);
        unimplemented!()
    }

    async fn download(&self, name: &String) -> Result<(), Box<dyn Error>> {
        if name != "default" {
            return Err(
                "only 'default' is currently supported: spc workspace download default".into(),
            );
        }

        let url = DEFAULT_WORKSPACE_URL;
        println!("Getting: {:?}", url);
        println!(
            "{} {}Validating url...",
            "[1/3]".colorize("bold yellow"),
            output::LOOKING_GLASS
        );

        let path = match parser::parse_url(url) {
            Ok(path) => path,
            Err(err) => {
                eprintln!("{}", err.to_string().colorize("red"));
                process::exit(0);
            }
        };

        println!("Default config URL: {}", path);

        let default_path: String = self
            .default_workspace_path
            .clone()
            .into_os_string()
            .into_string()
            .unwrap();

        println!("locating in dir: {}", default_path);
        let tmppath: PathBuf = PathBuf::from(default_path.clone());
        if !tmppath.exists() {
            println!("does not exist!");
            fs::create_dir_all(tmppath)?;
        }

        let data = match parser::parse_path(&path, Some(default_path)) {
            Ok(data) => data,
            Err(err) => {
                eprintln!("{}", err.to_string().colorize("red"));
                process::exit(0);
            }
        };

        println!(
            "{} {}Downloading...",
            "[2/3]".colorize("bold yellow"),
            output::TRUCK
        );

        match requests::fetch_data(&data).await {
            Err(err) => {
                eprintln!("{}", err.to_string().colorize("red"));
                process::exit(0);
            }
            Ok(_) => println!(
                "\n{} {}Downloaded Successfully.",
                "[3/3]".colorize("bold yellow"),
                output::SPARKLES
            ),
        };

        // if args.zipped {
        //     let dst_zip = format!("{}.zip", &data.root);
        //     let zipper = ZipArchiver::new(&data.root, &dst_zip);
        //     match zipper.run() {
        //         Ok(_) => (),
        //         Err(ZipError::FileNotFound) => {
        //             eprintln!(
        //                 "{}",
        //                 "\ncould not zip the downloaded file".colorize("bold red")
        //             )
        //         }
        //         Err(e) => eprintln!("{}", e.to_string().colorize("bold red")),
        //     }
        // }

        Ok(())
    }

    #[tokio::main]
    pub async fn manage(
        &self,
        cmd: String,
        workspace_name: String,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let workspace_path: PathBuf = self.default_workspace_path.clone();
        let target_workspace_path = self.default_workspace_path.join(&workspace_name);

        match cmd.as_str() {
            "download" => self.download(&workspace_name).await,
            "new" => self.cp(workspace_path, target_workspace_path),
            "rm" => self.rm(target_workspace_path),
            _ => Ok(()),
        }
    }
}
