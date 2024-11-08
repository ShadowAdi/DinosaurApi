import requests
from bs4 import BeautifulSoup

url = "https://www.dino-data.ca/browsedino.php"
r = requests.get(url)

# Check if the request was successful
if r.status_code == 200:
    soup = BeautifulSoup(r.text, "html.parser")
    dinos = []

    # Find all sections
    sections = soup.find_all("section")

    for section in sections:
        ul_tags = section.find_all("ul")
        for ul in ul_tags:
            li_tags = ul.find_all("li")
            for li in li_tags:
                a_tags = li.find_all("a")
                for a_tag in a_tags:
                    link = a_tag.get("href")
                    if link:
                        # Add base URL only if link is relative
                        full_link = "https://www.dino-data.ca/" + link if not link.startswith("http") else link
                        dinos.append(full_link)

    # Write links to file, each on a new line
    with open("links.txt", "w") as file:
        for link in dinos:
            file.write(link + "\n")

    print("Links saved successfully to links.txt.")
else:
    print(f"Failed to retrieve the page. Status code: {r.status_code}")
