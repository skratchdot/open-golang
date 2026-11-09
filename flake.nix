{
  description = "open-golang development environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          # The tests launch firefox via xdg-open. Provide a headless passthrough for curl.
          firefoxStub = pkgs.writeShellScriptBin "firefox" ''
            #!/usr/bin/env bash
            log_dir="$(dirname $(go env GOMOD))/logs"
            mkdir -p "''${log_dir}"
            log="$(mktemp -p "''${log_dir}" --suffix=".log" "$(date +%s).XXXXXXXXXX")"
            exec curl -v "''${@:1}" > "''${log}" 2>&1
          '';
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              curl
              delve
              go
              gopls
              golangci-lint
              xdg-utils
            ];
            CGO_ENABLED=0;
            XDG_CURRENT_DESKTOP="X-Generic";
            shellHook = ''
              export PATH="${firefoxStub}/bin:$PATH"
              export BROWSER=firefox
              # Keep the vendored xdg-open on its generic path so the host's
              # desktop session doesn't change test behavior.
              unset \
                  DBUS_SESSION_BUS_ADDRESS \
                  DESKTOP_SESSION \
                  DISPLAY \
                  GNOME_DESKTOP_SESSION_ID \
                  KDE_FULL_SESSION \
                  WAYLAND_DISPLAY
            '';
          };
        }
      );

      checks = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          # The tests launch firefox via xdg-open. Provide a headless passthrough for curl.
          firefoxStub = pkgs.writeShellScriptBin "firefox" ''
            #!/usr/bin/env bash
            log_dir="$(dirname $(go env GOMOD))/logs"
            mkdir -p "''${log_dir}"
            log="$(mktemp -p "''${log_dir}" --suffix=".log" "$(date +%s).XXXXXXXXXX")"
            exec curl -v "''${@:1}" > "''${log}" 2>&1
          '';

          mkGoCheck =
            name: nativeBuildInputs: command:
            pkgs.buildGoModule {
              pname = name;
              version = "0.1.0";
              src = ./.;
              inherit nativeBuildInputs;
              vendorHash = "sha256-qpdleQNbcjU1ZJoEAdRcMjXGXSVjMtsk7i/CyS/VVAg=";
              env.CGO_ENABLED = 0;
              env.XDG_CURRENT_DESKTOP="X-Generic";
              doCheck = true;
              checkPhase = ''
                export PATH="${firefoxStub}/bin:$PATH"
                export BROWSER=firefox
                unset \
                  DBUS_SESSION_BUS_ADDRESS \
                  DESKTOP_SESSION \
                  DISPLAY \
                  GNOME_DESKTOP_SESSION_ID \
                  KDE_FULL_SESSION \
                  WAYLAND_DISPLAY
                ${command}
              '';
            };
        in
        {
          go-test = mkGoCheck "go-test" [
            pkgs.curl
            pkgs.go
            pkgs.xdg-utils
            pkgs.writableTmpDirAsHomeHook
          ] "go test -v ./...";

          go-lint = mkGoCheck "go-lint" [
            pkgs.curl
            pkgs.go
            pkgs.golangci-lint
            pkgs.xdg-utils
            pkgs.writableTmpDirAsHomeHook
          ] "golangci-lint run ./...";
        }
      );
    };
}
