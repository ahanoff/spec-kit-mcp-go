

# Servidor MCP Spec-Kit (Go)

Un servidor del Protocolo de Contexto de Modelo (MCP) escrito en Go que conecta Amazon Q Developer con Spec-Kit de Microsoft, permitiendo flujos de trabajo de desarrollo basados en especificaciones a través de interacción por lenguaje natural.

## Características

- **init_project**: Inicializar nuevos proyectos spec-kit
- **specify**: Crear especificaciones de características a partir de lenguaje natural
- **plan**: Generar planes de implementación a partir de especificaciones
- **implement**: Generar código a partir de especificaciones
- **analyze**: Analizar el estado del proyecto y proporcionar información útil
- **tasks**: Desglosar el trabajo en tareas accionables

## Requisitos previos

- Go 1.21 o posterior
- Python 3.8 o posterior
- uv (gestor de paquetes de Python)
- CLI de Amazon Q Developer
- Spec-Kit instalado globalmente: `uv tool install specify-cli --from git+https://github.com/github/spec-kit.git`

## Instalación

1. Clonar este repositorio:
```bash
git clone git@github.com:ahanoff/spec-kit-mcp-go.git
cd spec-kit-mcp-go
```

2. Compilar el servidor MCP:
```bash
chmod +x build.sh
./build.sh
```

3. Instalar spec-kit globalmente:
```bash
uv tool install specify-cli --from git+https://github.com/github/spec-kit.git
```

## Configuración

Agregue el servidor MCP a su configuración de Q Developer en `~/.aws/amazonq/mcp.json`:

```json
{
  "mcpServers": {
    "spec-kit": {
      "type": "stdio",
      "command": "/path/to/spec-kit-mcp-go/spec-kit-mcp-go",
      "env": {
        "SPEC_KIT_WORKING_DIR": "/path/to/your/projects"
      },
      "timeout": 120000
    }
  }
}
```

## Uso

Inicie Q Developer y use lenguaje natural para interactuar con spec-kit:

```
Use la herramienta init_project para crear un nuevo proyecto spec-kit llamado "my-app" con Claude como asistente de IA
```

```
Use la herramienta specify para crear una especificación para la autenticación de usuario con correo electrónico y contraseña
```

## Desarrollo

Para modificar el servidor:

1. Realice sus cambios en `main.go`
2. Ejecute `go mod tidy` para actualizar las dependencias
3. Vuelva a compilar con `./build.sh`
4. Reinicie Q Developer

## Solución de problemas

- Asegúrese de que la ruta del binario en `mcp.json` sea correcta
- Verifique que spec-kit esté instalado: `specify --version`
- Verifique que Go esté instalado: `go version`
- Pruebe el servidor manualmente: `echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./spec-kit-mcp-go`

## Licencia

MIT
