"""
    This rule is replacing constant variables with their values
    For example:
        x = 10; print(x + 100)
    Transformed into:
        print(10 + 100)
"""
import ast

from weasel.const_counter import ConstCounter


class Rule_1(ast.NodeTransformer):
    def __init__(self, statistic, root=None):
        super().__init__()
        self._root = root
        self._statistic = statistic

    def visit_FunctionDef(self, node):
        if node != self._root:
            locals = ConstCounter()
            for child in node.body:
                locals.visit(child)
            return Rule_1(self._statistic + [locals], root=node).visit(node)
        new_body = []
        for child in node.body:
            new_child = self.visit(child)
            if new_child is None:
                continue
            new_body.append(ast.copy_location(new_child, child))
        # If we leave function body empty it's runtime error
        if len(new_body) == 0:
            new_body.append(ast.copy_location(ast.Pass(), node.body[0]))
        return ast.copy_location(
            ast.FunctionDef(
                name=node.name,
                args=node.args,
                body=new_body,
                decorator_list=node.decorator_list,
                returns=node.returns,
            ),
            node,
        )

    def visit_Name(self, node):
        # Note: order is important
        for lookup in self._statistic[::-1]:
            if node.id in lookup._statistic:
                v = lookup._statistic[node.id]
                if v.count > 1:
                    return node
                return ast.copy_location(
                    ast.Constant(
                        value=v.value.value,
                        kind=v.value.kind,
                        s=v.value.s,
                        n=v.value.n,
                    ),
                    node,
                )
        return node

    def visit_Assign(self, node):
        # TODO: we want to optimize such cases also
        if len(node.targets) > 1 or not isinstance(node.targets[0], ast.Name):
            return node
        target_name = node.targets[0].id
        if isinstance(node.value, ast.Constant):
            for lookup in self._statistic[::-1]:
                if target_name in lookup._statistic:
                    v = lookup._statistic[target_name]
                    if v.count > 1:
                        return node
                    return None
        # Visiting individual expression nodes
        # a = b + c
        #     ^ ^ ^
        for child in ast.iter_child_nodes(node):
            self.visit(child)
        return node
