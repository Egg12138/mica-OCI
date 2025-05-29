pub mod create;
pub mod delete;

#[derive(Parser, Debug)]
pub enum MinimalCmd {
    Create(Create),
    Start(Start),
    State(State),
    Kill(Kill),
    Delete(Delete),
}

#[derive(Parser, Debug)]
pub enum CommonCmd {
    Checkpointt(Checkpoint),
    Events(Events),
    Exec(Exec),
    Features(Features),
    List(List),
    Pause(Pause),
    #[clap(allow_hyphen_values = true)]
    Ps(Ps),
    Resume(Resume),
    Run(Run),
    Update(Update),
    Spec(Spec),
}