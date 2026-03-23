{
  description = "Golang flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: 
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
  in 
  {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = with pkgs; [
        starship
        go
        gopls
        gotools
        just
      ];
      # shellHook sets the SHELL envar because there is a bug in Nixos
        # that incorrectly loads a sub-shell, causing broken prompts.
        # It should have been fixed int 25.11 but I am still running into it 
        # see: https://github.com/NixOS/nix/issues/5873
      shellHook = ''
        export SHELL=${ pkgs.lib.getExe pkgs.bash }
        export PATH=$PATH:$GOPATH/bin
        eval "$(starship init bash)"
        export STARSHIP_CONFIG=$PWD/starship.toml
      '';
    };
      
  };
}
