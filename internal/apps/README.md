# App Platform MCP Tools

This directory contains tools for the App Platform MCP Server. These tools are used to interact with the App Platform API and perform various operations on App Platform resources.

## Supported Tools

- `create-app-from-spec`: Initialize a new App Platform app by connecting a GitHub, GitLab, or Bitbucket repo (including specifying the branch and build settings). Condenses the app creation workflow into one action for the agent.
- `apps-update`: Modify an app’s settings or trigger a re-deploy. Allows updating environment variables, scaling parameters, redeploying, etc.
- `apps-delete`: Delete an App Platform app.
- `apps-get-info`: Get the details and status of an existing app.
- `apps-usage`: Get live information about an app’s resource usage, like CPU and memory consumption.
- `apps-get-deployment-status`: Check the status of a specific deployment for an App Platform app.
- `apps-list`: List all App Platform apps in the account.

# Example queries using App Platform MCP Tools

- Can you deploy this app from this git repository?
- Show me all of my apps in app platform.
- Delete this application for me.
- Give me the deployment status of this app.
- Which environment variables are set for this app?
- Trigger a new deployment for my app.
- Update the instance size for my app.