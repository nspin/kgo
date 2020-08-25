with import <nixpkgs> {};

mkShell {
  nativeBuildInputs = [ go_1_14 ];
  shellHook = ''
    GOPATH=$HOME/go
  '';
}
