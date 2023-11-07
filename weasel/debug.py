def print_node(node, indent):
    print(" " * indent, node, sep="")
    for child in getattr(node, "body", []):
        print_node(child, indent + 4)
