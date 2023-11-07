"""
    This rule is replacing constant variables with their values
    For example:
        x = 10
        print(x + 100)
    Transformed into:
        print(10 + 100)
"""
import ast
import copy

# TODO: add testcases


class Rule_1(ast.NodeTransformer):
    def __init__(self, root=None, globals=None):
        super().__init__()
        self._root = root
        self._global_constants = globals or {}
        self._local_constants = {}

    def visit_FunctionDef(self, node):
        if node != self._root:
            new_globals = copy.deepcopy(self._global_constants)
            new_globals.update(self._local_constants)
            return Rule_1(root=node, globals=new_globals).visit(node)
        new_body = []
        for child in node.body:
            new_child = self.visit(child)
            if new_child is None:
                continue
            new_body.append(ast.copy_location(new_child, child))
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
        for lookup_dict in (self._local_constants, self._global_constants):
            if node.id in lookup_dict:
                const = lookup_dict[node.id]
                return ast.copy_location(
                    ast.Constant(
                        value=const.value,
                        kind=const.kind,
                        s=const.s,
                        n=const.n,
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
            self._local_constants[target_name] = node.value
            return None
        elif target_name in self._local_constants:
            self._local_constants.pop(target_name)
            self._global_constants.pop(target_name, None)
        elif target_name in self._global_constants:
            self._global_constants.pop(target_name, None)
        # Visiting individual expression nodes
        # a = b + c
        #     ^ ^ ^
        for child in ast.iter_child_nodes(node):
            self.visit(child)
        return node
