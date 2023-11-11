import ast

from weasel.transformers.rule_1 import Rule_1, Rule_1_visitor


def optimize(root: ast.AST) -> ast.AST:
    visitor = Rule_1_visitor()
    visitor.visit(root)
    root = Rule_1([visitor]).visit(root)
    return root
