__author__ = 'creed'

# #import sys
# #import os
# sys.path.append(os.path.abspath(os.path.join('../../golem/apps/core/benchmark', 'golem')))
#
# # from benchmarkrunner import BenchmarkRunner
# # from apps.core.benchmark.benchmarkrunner import BenchmarkRunner
# # from golem.apps.core.benchmark import BenchmarkRunner

def timed_system_test(command):
    import time
    start = time.time()
    os.system(command)
    return time.time() - start

def timed_test( scene_file, resolution ):
    command = 'blender '
    command += '-noaudio '
    command += '-b ' + scene_name + ' '
    command += '--python-expr '
    command += '"import bpy; bpy.data.scenes[ 0 ].render.resolution_x = ' + x + '; bpy.data.scenes[ 0 ].render.resolution_y = ' + y + ';" '
    return timed_system_test(command)

def estimate_time( scene_file, resolution ):
    s_1 = 0.1
    s_2 = 0.05

    res_1 = resolution * s_1
    res_2 = resolution * s_2

    t_1 = timed_test( scene_file, res_1)
    t_2 = timed_test( scene_file, res_2)

    return t_1 * ( s_2 ** 2 - 1 ) / ( s_2 ** 2 - s_1 ** 2 ) +\
            t_2 * ( 1 - s_1 ** 2 ) / ( s_2 ** 2 - s_1 ** 2 )
