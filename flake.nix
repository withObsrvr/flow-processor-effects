{
  description = "Obsrvr Flow Plugin: Effects Processor";

  nixConfig = {
    allow-dirty = true;
  };

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages = {
          default = pkgs.buildGoModule {
            pname = "flow-processor-effects";
            version = "0.1.0";
            src = ./.;
            vendorHash = null;
            hardeningDisable = [ "all" ];
            preBuild = ''
              export CGO_ENABLED=1
            '';
            buildPhase = ''
              runHook preBuild
              go build -buildmode=plugin -o flow-processor-effects.so .
              runHook postBuild
            '';
            installPhase = ''
              runHook preInstall
              mkdir -p $out/lib
              cp flow-processor-effects.so $out/lib/
              mkdir -p $out/share
              cp go.mod $out/share/
              if [ -f go.sum ]; then
                cp go.sum $out/share/
              fi
              runHook postInstall
            '';
            nativeBuildInputs = [ pkgs.pkg-config ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_21
            pkg-config
            git
            gopls
            delve
          ];
          shellHook = ''
            export CGO_ENABLED=1
            if [ ! -d vendor ]; then
              echo "Vendoring dependencies..."
              go mod tidy
              go mod vendor
            fi
            echo "Development environment ready! To build, run: go build -buildmode=plugin -o flow-processor-effects.so ."
          '';
        };
      }
    );
} 