{
  description = "Flake to work in this project";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    
  };

  outputs = { self, nixpkgs }: 
    let 
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      
      nixpkgsFor = forAllSystems (system : import nixpkgs {
        inherit system; 
        config.allowUnfree = true;
        });
      
    in 
    {
      devShells = forAllSystems ( system: 
        let 
          pkgs = nixpkgsFor.${system};
        in 
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [go gopls gotools go-tools gnumake graphviz go-task delve gdlv];
          };
        });
    };
}