import os
import sys
import platform
import socket
import psutil
import subprocess
import uuid
import time
from datetime import datetime
import requests
import json

"""
This is in active development and is NOT ready for production use.

I am working to reverse engineer the endar agent, which may take some time. 
Please be patient.

Steps:
- Pull existing policies, if receiving 424 response, register agent
- register agent
- handle data collecxtion (performance, disk, etc)
"""

class SystemInfoCollector:
    def __init__(self, registration_token):
        self.registration_token = registration_token
    
    def get_windows_domain_info(self):
        """ Retrieve domain-related information on Windows """
        try:
            import wmi
            c = wmi.WMI()
            domain_info = {
                "domain": None,
                "forest": None,
                "dn": None,
                "site": None,
                "domain_joined": False,
                "is_dc": False
            }
            for obj in c.Win32_ComputerSystem():
                domain_info["domain"] = obj.Domain
                domain_info["domain_joined"] = obj.PartOfDomain
                domain_info["is_dc"] = "Primary Domain Controller" in obj.Roles if obj.Roles else False
            for obj in c.Win32_NTDomain():
                domain_info["forest"] = obj.DnsForestName
            dn_output = subprocess.run(["dsquery", "computer", f"name={socket.gethostname()}"] , capture_output=True, text=True, shell=True)
            domain_info["dn"] = dn_output.stdout.strip().split("\n")[0] if dn_output.returncode == 0 else "Unknown"
            site_output = subprocess.run(["nltest", "/dsgetsite"], capture_output=True, text=True, shell=True)
            domain_info["site"] = site_output.stdout.strip().split("\n")[0] if site_output.returncode == 0 else "Unknown"
            return domain_info
        except Exception:
            return {}

    def collect_system_data(self):
        """
        Collect system data from the user's machine.
        """
        hostname = socket.gethostname()
        fqdn = socket.getfqdn()
        local_ip = socket.gethostbyname(hostname)
        system = platform.system()
        release = platform.release()
        version = platform.version()
        machine = platform.machine()
        processor = platform.processor()
        cpu_count = psutil.cpu_count(logical=False)
        logical_cpu_count = psutil.cpu_count(logical=True)
        memory = f"{round(psutil.virtual_memory().total / (1024 ** 3))} GB"
        uptime = time.time() - psutil.boot_time()
        last_boot = datetime.fromtimestamp(psutil.boot_time()).strftime("%Y-%m-%d %H:%M:%S")
        
        install_type = "full"
        edition = "Unknown"
        build = "Unknown"
        domain_info = {
            "domain": None,
            "forest": None,
            "dn": None,
            "site": None,
            "domain_joined": False,
            "is_dc": False
        }
        
        if system == "Windows":
            import winreg
            try:
                with winreg.OpenKey(winreg.HKEY_LOCAL_MACHINE, r"SOFTWARE\Microsoft\Windows NT\CurrentVersion") as key:
                    edition = winreg.QueryValueEx(key, "ProductName")[0]
                    build = winreg.QueryValueEx(key, "CurrentBuild")[0]
            except Exception:
                pass
            print("Getting Windows domain info...")
            domain_info = self.get_windows_domain_info()
            print(domain_info)
        
        data = {
            "key": hostname,
            "token": self.registration_token,
            "version": "1.0.0",
            "enabled": True,
            "install_group": "default",
            "public_addr": local_ip,
            "country_code": "US",  # Needs external API call to get accurate data
            "country_name": "United States",  # Placeholder
            "region_name": "California",  # Placeholder
            "city_name": "San Francisco",  # Placeholder
            "lat": "37.7749",  # Placeholder
            "long": "-122.4194",  # Placeholder
            "uninstall": False,
            "cpu_count": cpu_count,
            "logical_cpu_count": logical_cpu_count,
            "memory": memory,
            "hostname": hostname,
            "fqdn": fqdn,
            "domain": domain_info["domain"],
            "forest": domain_info["forest"],
            "dn": domain_info["dn"],
            "site": domain_info["site"],
            "domain_joined": domain_info["domain_joined"],
            "is_dc": domain_info["is_dc"],
            "family": system,
            "release": release,
            "sys_version": version,
            "install_type": install_type,
            "edition": edition,
            "build": build,
            "machine": machine,
            "local_addr": local_ip,
            "processor": processor,
            "last_boot": last_boot,
            "last_active": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            "svc_start": last_boot,
            "svc_uptime": int(uptime),
            "uptime": int(uptime)
        }

        return data

class Agent:
    def __init__(self, server_url, registration_token):
        self.server_url = server_url
        self.registration_token = registration_token
        self.agent_id = None

        self.system_info_collector = SystemInfoCollector(registration_token)
        self.agent_info = {}

    def register_agent(self):
        """
        Register the agent with the server.
        """
        self.agent_info = self.system_info_collector.collect_system_data()
        print(self.agent_info)

        headers = {
            "Content-Type": "application/json",
            "tenant-key": self.registration_token
        }

        try:
            response = requests.post(f'{self.server_url}/api/v1/agent/register', headers=headers, json=self.agent_info)

            print("Status Code:", response.status_code)
            print("Response Headers:", response.headers)
            print("Response Body:", response.text)

            # Check if the response contains JSON
            if response.headers.get("Content-Type", "").startswith("application/json"):
                response_data = response.json()
                print(f"Response From Server: {response_data}")
                if response.status_code == 200 and response_data.get("registered"):
                    self.agent_id = response_data.get("agent_id")
                    print(f"Agent registered with ID: {self.agent_id}")
                else:
                    print(f"Error: {response_data.get('msg', 'Unknown error')}")
            else:
                print("Received non-JSON response from server:")
                print(response.text)

        except requests.exceptions.RequestException as e:
            print(f"Error connecting to server: {e}")

        return self.agent_id
    
def main():
    server_url = 'http://192.168.4.142'  # replace with your server URL
    registration_token = 'd16757ade4c6489290a09749f469f1af'  # replace with your registration token
    agent = Agent(server_url, registration_token)
    agent_id = agent.register_agent()

if __name__ == '__main__':
    main()