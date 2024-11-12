import io
import requests
import zipfile
import re
import os

token = os.getenv("GITHUB_TOKEN")
headers = {"Authorization": f"token {token}"}

res = {}

def getTPS(url):
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        zip_file = zipfile.ZipFile(io.BytesIO(response.content))

        for file_name in zip_file.namelist():
            if 'Build Chain' not in file_name: continue
            with zip_file.open(file_name) as file:
                content = file.read().decode('utf-8')
                lastLine = content.splitlines()[-1]
                print(lastLine)
                pattern = r"Best TPS: (\d+) GasUsed%: ([\d.]+)"

                match = re.search(pattern, lastLine)

                if match:
                    best_tps = int(match.group(1))
                    gas_used = float(match.group(2))
                    
                    return best_tps, gas_used
    return None, None

resp = requests.get('https://api.github.com/repos/0glabs/evmchainbench/actions/workflows', headers=headers).json()

for workflow in resp['workflows']:
    name = workflow['name'].split('-')
    if len(name) < 2: continue
    chain = name[1].strip()
    category = 'Simple' if len(name) < 3 else name[2].strip()
    print(chain, category)
    runs = requests.get(f'{workflow['url']}/runs', headers=headers).json()
    last_run = runs['workflow_runs'][0]
    print(last_run['logs_url'])
    best_tps, gas_used = getTPS(last_run['logs_url'])
    print(best_tps, gas_used)
    if chain not in res:
        res[chain] = {}
    res[chain][category] = (best_tps, gas_used)

print(res)

from prettytable import PrettyTable

table = PrettyTable()
table.field_names = ["Chain", "Simple", "ERC20", "Uniswap"]

for chain, contracts in res.items():
    row = [chain]
    for contract_type in ["Simple", "ERC20", "Uniswap"]:
        tps, gas_used = contracts[contract_type]
        if tps is None or gas_used is None:
            row.append("")
        else:
            row.append(f"{tps:4}, {gas_used * 100:.2f}%")
    table.add_row(row)

markdown_table = "| " + " | ".join(table.field_names) + " |\n"
markdown_table += "| " + " | ".join(["---"] * len(table.field_names)) + " |\n"
for row in table.rows:
    markdown_table += "| " + " | ".join(map(str, row)) + " |\n"

with open(os.getenv('GITHUB_STEP_SUMMARY'), 'a') as summary_file:
    summary_file.write("### Performance Table\n")
    summary_file.write(markdown_table + "\n")
