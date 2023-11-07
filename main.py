import ast
import importlib

from weasel.debug import print_node
from weasel.lib import optimize


def main():
    root = ast.parse(open("example.py", "r").read())
    optimized = optimize(root)
    compiled = compile(optimized, "example.py", "exec")
    bytecode = importlib._bootstrap_external._code_to_hash_pyc(
        compiled,
        b"12345678",
        True,
    )
    with open("example.pyc", "wb") as fd:
        fd.write(bytecode)


if __name__ == "__main__":
    main()
