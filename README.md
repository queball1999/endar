<div align="center">
  <img src="/img/endar_rounded.png" height="150"/>
</div>

<h1 align="center" style="margin-top: -5px"> Endar </h1>
<p align="center" style="width: 100;">
   <a href="https://github.com/queball1999/endar/wiki">Docs</a>
   |
   <a href="https://www.postman.com/qssapi/endar/overview">Postman Workspace</a>
   |
   <a href="https://discord.gg/TkSKvCk6rf">Discord</a>
</p>

This is a continuation of [tomkeene's](https://github.com/tomkeene) project, [endar](https://github.com/tomkeene/endar).

### Table of Contents
1. [About](#about)
2. [Main Features](#main-features)
3. [Getting Started](#getting-started)
4. [Creating your first policy](#creating-your-first-policy)
5. [Roadmap](#roadmap)

:bulb: :zap: :bulb: [View the video of Endar here!](https://drive.google.com/file/d/1CJolj-nP7z19-5DtQQRwgZvGTq9Ej75c/view)

### About
Endar is an RMM (Remote monitoring and management) tool that supports Windows, Linux and MacOS. The Endar agent is a binary that runs on the endpoints and communicates with the Endar server. *Endar is currently in Alpha mode - while it works great, the server architecture does not support 100's of agents.*


Home Dashboard          |
:-------------------------:|
![](img/endar_dash.PNG)  |

### Main Features
Endar has two main features that are simple; Compliance Management & Monitoring. This tool was originially created to meet compliance requirements.

##### Compliance Management
Compliance management allows you to ask if something is true (assertion) and then optionally perform remediation. In practice, Endar uses scripts/programs to achieve this. For example, lets pretend you want to ensure the Windows firewall is enabled (a common compliance task). You would create a "validation" script to check if the firewall is enabled. If the firewall is _not_ enabled, your "enforcement" script would then execute, bringing the asset back into compliance. Endar is not opinionated so you can write scripts in whatever language you please (powershell, bash, python, etc).

Compliance Management          |  
:-------------------------:|
![](img/endar_comp.PNG)  |

##### Monitoring
Monitoring consists of the Endar agents collecting host-based metrics and periodically forwarding the data to the Endar server. Metrics consist of disk performance, memory stats, disk stats, load performance and more.

Monitoring Agents          |  
:-------------------------:|
![](img/endar_perf.PNG)  |

### Getting Started

To get started with Endar, visit our dedicated [Getting Started Guide](https://github.com/queball1999/endar/wiki/Getting-Started) in the Wiki.

This guide covers:
- Setting up the server with Docker
- Configuring and deploying agents for Windows, Linux, and MacOS
- Creating your first compliance policy
- And more!

If you have questions or run into issues, please refer to the guide for detailed instructions or open an issue on GitHub.

### Roadmap
For the latest project plans and progress, please visit our updated [roadmap](https://github.com/queball1999/endar/projects?query=is%3Aopen).

Below is a list of improvements desired by the original developer. We will do our best to review and possibly implement the following:
- [ ] Improve monitoring to gather software, services, scheduled tasks (cronjobs), users and groups
- [ ] Improve monitoring to collect process specific metrics
- [ ] Improve deployment - currently a binary is provided by OS specific installers would be nice
- [ ] Improve architecture - the current deployment architecture will not support hundreds of agents. Its probably easiest to just leverage managed solutions of a popular provider such as GCP.  
- [ ] Add other "RMM" features (one-click install of apps, etc)
- [ ] Upload source code of clients to repo
