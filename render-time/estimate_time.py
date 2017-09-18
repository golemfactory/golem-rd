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
    import os
    start = time.time()
    os.system(command)
    return time.time() - start

def timed_test( scene_file, resolution ):
    x = resolution[ 0 ]
    y = resolution[ 1 ]

    command = 'blender '
    command += '-noaudio '
    command += '-b ' + scene_file + ' '
    command += '--python-expr '
    command += '"import bpy; bpy.data.scenes[ 0 ].render.resolution_x = ' + str( x ) + '; bpy.data.scenes[ 0 ].render.resolution_y = ' + str( y ) + '; '
    command += 'bpy.data.scenes[ 0 ].render.resolution_percentage=100;" '
    command += '-f 1 '
    print 'command: ' + command
    return timed_system_test(command)

def scale_res( resolution, scale ):
    return [ r * scale for r in resolution ]

def estimate_time( scene_file, resolution ):
    s_1 = 0.05
    s_2 = 0.1

    res_1 = scale_res( resolution, s_1 )
    res_2 = scale_res( resolution , s_2 )

    t_1 = timed_test( scene_file, res_1)
    t_2 = timed_test( scene_file, res_2)

    return t_1 * ( s_2 ** 2 - 1 ) / ( s_2 ** 2 - s_1 ** 2 ) +\
            t_2 * ( 1 - s_1 ** 2 ) / ( s_2 ** 2 - s_1 ** 2 )

if __name__ == "__main__":
    print estimate_time( 'bmw27_cpu.blend', [1000.0 ,1000.0 ] )