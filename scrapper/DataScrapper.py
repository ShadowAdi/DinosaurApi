import requests
from bs4 import BeautifulSoup
import json

# Read URLs from the file
with open("links.txt", "r") as file:
    lines = file.readlines()
url = "https://www.dino-data.ca/browsedino.php"

data = []  # To store data for all dinosaurs

for line in lines:
    url = line.strip()  # Remove any newline characters
    r = requests.get(url)

    # Check if the request was successful
    if r.status_code == 200:
        soup = BeautifulSoup(r.text, "html.parser")

        # Initialize a dictionary to store information for each dinosaur
        dino = {}

        InfoContent = soup.find("div", id="InfoContent")
        if InfoContent:
            # Image
            InfoImage = InfoContent.find("img", id="InfoImage")
            if InfoImage and InfoImage.has_attr("src"):
                dino["img"] = "https://www.dino-data.ca/" + InfoImage["src"]
            else:
                dino["img"] = None  # or an alternative value

            # Dino name
            DinoName = InfoContent.find("h1", id="DinoName")
            dino["dinoName"] = DinoName.text if DinoName else "No Name"

            # Small description
            DinoSubdescription = InfoContent.find("p", class_="dino-subtitle")
            dino["dinoSmallDescription"] = (
                DinoSubdescription.text if DinoSubdescription else "No Description"
            )

            # Full description
            article = InfoContent.find("article", id="DinoWriteUp")
            if article:
                p_tag = article.find("p")
                dino["description"] = (
                    p_tag.get_text(strip=True) if p_tag else "No description available"
                )
            else:
                dino["description"] = "No article available"

            # Paleontologists
            design_div = InfoContent.find("div", class_="w3-border-top w3-center")
            if design_div:
                paleontologists_p = design_div.find_all("p", class_="w3-small")

                # Make sure we are targeting the second <p> tag with `w3-small`
                if len(paleontologists_p) > 1:
                    paleontologists_links = paleontologists_p[1].find_all("a")

                    # Extract href and text for each <a> and prepend base URL
                    dino["paleontologists"] = [
                        {"name": a.text, "link": url + a["href"]}
                        for a in paleontologists_links
                        if a.has_attr("href")
                    ]
                else:
                    dino["paleontologists"] = []
            else:
                dino["paleontologists"] = []

            # Extract additional details from 'w3-cell-row'
            W3cellsrows = InfoContent.find("section", class_="w3-cell-row")
            if W3cellsrows:
                W3cells = W3cellsrows.find_all("div", class_="w3-cell")
                for i, w3cell in enumerate(W3cells):
                    p_tags = w3cell.find_all("p")
                    for j, p_tag in enumerate(p_tags):
                        b_tag = p_tag.find("b")
                        if b_tag:
                            key = b_tag.get_text(strip=True)
                            value = p_tag.get_text(strip=True).replace(key, "").strip()
                            dino[key] = value

                        # For the last <p> tag in the second cell, capture <a> tags as an array
                        if i == 1:
                            if j == len(p_tags) - 1:
                                dino[key] = [
                                    a.get_text(strip=True) for a in p_tag.find_all("a")
                                ]
            else:
                print("No W3cellsrows found.")

            # Append the dinosaur data to the list
            data.append(dino)
        else:
            print("No InfoContent id found for URL:", url)
    else:
        print(f"Failed to retrieve the page at {url}. Status code: {r.status_code}")


with open("Data.json", "w") as file:
    json.dump(data, file, indent=4)
