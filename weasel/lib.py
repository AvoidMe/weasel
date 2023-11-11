import ast

from weasel.const_counter import ConstCounter
from weasel.transformers.rule_1 import Rule_1


def optimize(root: ast.AST) -> ast.AST:
    visitor = ConstCounter()
    visitor.visit(root)
    root = Rule_1([visitor]).visit(root)
    return root
