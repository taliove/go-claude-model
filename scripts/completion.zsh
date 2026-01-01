#compdef ccm

autoload -U compinit
compinit

_ccm() {
  local -a commands
  commands=(
    'add:Add or configure a provider'
    'edit:Update provider configuration'
    'remove:Remove a configured provider'
    'list:List all providers'
    'show:Show provider details'
    'run:Run Claude Code with specified provider'
    'test:Test provider API connection'
    'generate:Generate shell scripts for providers'
    'switch:Interactively switch provider'
    'init:First-time setup wizard'
    'version:Show version information'
    'help:Show help'
  )

  _arguments -C \
    '1: :_ccm_commands' \
    '*::arg:->args'

  case $state in
    args)
      local cmd=$words[1]
      case $cmd in
        add|edit|remove|run|test|show)
          _arguments \
            '(--key -k)'{-k,--key}'[API key]' \
            '(--url -u)'{-u,--url}'[API URL]' \
            '(--model -m)'{-m,--model}'[Model name]' \
            '(--force -f)'{-f,--force}'[Force operation]'
          ;;
        generate|switch|init|version|list|help)
          ;;
      esac
      ;;
  esac
}

_ccm_commands() {
  _describe -t commands 'ccm command' commands
}

_ccm "$@"
