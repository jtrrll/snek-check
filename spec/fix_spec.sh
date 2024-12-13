# shellcheck shell=sh

Describe "fix"
  Context "with an empty, valid directory"
    It "succeeds"
      When call "$bin" --fix "$root"
      The status should be success
    End
  End

  Context "with a single invalid file"
    create_invalid_file() { touch "$root"/InVaLiD; }
    BeforeEach "create_invalid_file"

    It "succeeds"
      When call "$bin" --fix "$root"
      The status should be success
    End

    It "renames the file"
      When call "$bin" --fix "$root"
      The file "$root"/InVaLiD should not be exist
      The file "$root"/invalid should be exist
    End
  End
End
